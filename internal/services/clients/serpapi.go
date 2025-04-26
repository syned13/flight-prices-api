package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
)

type SerpAPIClient struct {
}

func NewSerpAPIClient() *SerpAPIClient {
	return &SerpAPIClient{}
}

func (c *SerpAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
