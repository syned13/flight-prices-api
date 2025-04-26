package itinerary_fetcher

import (
	"context"
	"sync"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/internal/services/clients"
)

type ItineraryFetcherService interface {
	FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error)
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

func (s *itineraryFetcherService) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
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
