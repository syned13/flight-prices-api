package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
)

type AmadeusClient struct {
}

func NewAmadeusClient() *AmadeusClient {
	return &AmadeusClient{}
}

func (c *AmadeusClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
