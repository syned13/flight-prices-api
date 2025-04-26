package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type AmadeusClient struct {
	baseURL string
	apiKey  string
}

func NewAmadeusClient() *AmadeusClient {
	return &AmadeusClient{
		baseURL: config.GetConfig().Amadeus().BaseURL(),
		apiKey:  config.GetConfig().Amadeus().APIKey(),
	}
}

func (c *AmadeusClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
