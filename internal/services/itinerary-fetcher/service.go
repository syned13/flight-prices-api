package itinerary_fetcher

import (
	"context"
	"errors"
	"log"
	"sort"
	"strconv"
	"sync"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/internal/services/clients"
)

type ItineraryFetcherService interface {
	FetchItineraries(ctx context.Context, request models.FlightSearchRequest) (*models.FlightSearchResponse, error)
}

type itineraryFetcherService struct {
	amadeusClient   clients.ItineraryFetcherClient
	flightAPIClient clients.ItineraryFetcherClient
	serpAPIClient   clients.ItineraryFetcherClient
}

func NewItineraryFetcherService() ItineraryFetcherService {
	clientFactory := clients.NewClientFactory()
	amadeusClient := clientFactory.NewItineraryFetcherClient(clients.Amadeus)
	flightAPIClient := clientFactory.NewItineraryFetcherClient(clients.FlightAPI)
	serpAPIClient := clientFactory.NewItineraryFetcherClient(clients.SerpAPI)

	return &itineraryFetcherService{
		amadeusClient:   amadeusClient,
		flightAPIClient: flightAPIClient,
		serpAPIClient:   serpAPIClient,
	}
}

func (s *itineraryFetcherService) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) (*models.FlightSearchResponse, error) {
	itineraries, err := s.fetchItineraries(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(itineraries) == 0 {
		return nil, errors.New("no itineraries found")
	}

	sortedByPrice := sortItinerariesByPrice(itineraries)
	sortedByDuration := sortItinerariesByDuration(itineraries)

	cheapest := sortedByPrice[0]
	fastest := sortedByDuration[0]

	return &models.FlightSearchResponse{
		Itineraries: itineraries,
		Cheapest:    cheapest,
		Fastest:     fastest,
	}, nil
}

func (s *itineraryFetcherService) fetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	type result struct {
		itineraries []models.Itinerary
		err         error
		source      string
	}

	resultChan := make(chan result, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		itineraries, err := s.amadeusClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err, source: "Amadeus"}
	}()

	go func() {
		defer wg.Done()
		itineraries, err := s.flightAPIClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err, source: "FlightAPI"}
	}()

	go func() {
		defer wg.Done()
		itineraries, err := s.serpAPIClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err, source: "SerpAPI"}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allItineraries []models.Itinerary
	failedFetches := 0

	for res := range resultChan {
		if res.err != nil {
			log.Printf("Error fetching itineraries from %s: %v", res.source, res.err)
			failedFetches++
			continue
		}

		log.Printf("Received %d itineraries from %s", len(res.itineraries), res.source)
		for _, itin := range res.itineraries {
			log.Printf("%s price: %s %s", res.source, itin.Price.Total, itin.Price.Currency)
		}

		allItineraries = append(allItineraries, res.itineraries...)
	}

	// Only return error if all fetches failed
	if failedFetches == 3 {
		return nil, errors.New("all flight price fetches failed")
	}

	if len(allItineraries) == 0 {
		return nil, errors.New("no itineraries found")
	}

	log.Printf("All prices before sorting:")
	for _, itin := range allItineraries {
		log.Printf("Price: %s %s", itin.Price.Total, itin.Price.Currency)
	}

	return allItineraries, nil
}

func sortItinerariesByPrice(itineraries []models.Itinerary) []models.Itinerary {
	sorted := make([]models.Itinerary, len(itineraries))
	copy(sorted, itineraries)

	sort.Slice(sorted, func(i, j int) bool {
		priceI, errI := strconv.ParseFloat(sorted[i].Price.Total, 64)
		priceJ, errJ := strconv.ParseFloat(sorted[j].Price.Total, 64)

		if errI != nil {
			log.Printf("Error parsing price %s: %v", sorted[i].Price.Total, errI)
			return false
		}
		if errJ != nil {
			log.Printf("Error parsing price %s: %v", sorted[j].Price.Total, errJ)
			return true
		}

		return priceI < priceJ
	})

	log.Printf("Sorted prices:")
	for _, itin := range sorted {
		log.Printf("Price: %s %s", itin.Price.Total, itin.Price.Currency)
	}

	return sorted
}

func sortItinerariesByDuration(itineraries []models.Itinerary) []models.Itinerary {
	sort.Slice(itineraries, func(i, j int) bool {
		return itineraries[i].DurationInMinutes < itineraries[j].DurationInMinutes
	})
	return itineraries
}
