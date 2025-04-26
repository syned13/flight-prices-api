package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type FlightAPIClient struct {
	baseURL string
	apiKey  string
}

func NewFlightAPIClient() *FlightAPIClient {
	return &FlightAPIClient{
		baseURL: config.GetConfig().FlightAPI().BaseURL(),
		apiKey:  config.GetConfig().FlightAPI().APIKey(),
	}
}

func (c *FlightAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
