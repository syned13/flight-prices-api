package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
)

const (
	Amadeus   = "amadeus"
	FlightAPI = "flightapi"
	SerpAPI   = "serpapi"
)

type ClientFactory struct {
}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{}
}

type ItineraryFetcherClient interface {
	FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error)
}

func (f *ClientFactory) NewItineraryFetcherClient(clientType string) ItineraryFetcherClient {
	switch clientType {
	case Amadeus:
		return NewAmadeusClient()
	case FlightAPI:
		return NewFlightAPIClient()
	case SerpAPI:
		return NewSerpAPIClient()
	default:
		return nil
	}
}
