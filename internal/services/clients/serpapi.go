package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type SerpAPIClient struct {
	baseURL string
	apiKey  string
}

func NewSerpAPIClient() *SerpAPIClient {
	return &SerpAPIClient{
		baseURL: config.GetConfig().SerpAPI().BaseURL(),
		apiKey:  config.GetConfig().SerpAPI().APIKey(),
	}
}

func (c *SerpAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
