package itinerary_fetcher

import (
	"context"
	"errors"
	"sort"
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
	}

	resultChan := make(chan result, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		itineraries, err := s.amadeusClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err}
	}()

	go func() {
		defer wg.Done()
		itineraries, err := s.flightAPIClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err}
	}()

	go func() {
		defer wg.Done()
		itineraries, err := s.serpAPIClient.FetchItineraries(ctx, request)
		resultChan <- result{itineraries: itineraries, err: err}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allItineraries []models.Itinerary
	for res := range resultChan {
		if res.err != nil {
			return nil, res.err
		}
		allItineraries = append(allItineraries, res.itineraries...)
	}

	return allItineraries, nil
}

func sortItinerariesByPrice(itineraries []models.Itinerary) []models.Itinerary {
	sort.Slice(itineraries, func(i, j int) bool {
		return itineraries[i].Price.Total < itineraries[j].Price.Total
	})
	return itineraries
}

func sortItinerariesByDuration(itineraries []models.Itinerary) []models.Itinerary {
	sort.Slice(itineraries, func(i, j int) bool {
		return itineraries[i].Duration < itineraries[j].Duration
	})

	return itineraries
}
