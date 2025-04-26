package clients

import (
	"context"

	"github.com/syned13/flight-prices-api/internal/models"
)

type FlightAPIClient struct {
}

func NewFlightAPIClient() *FlightAPIClient {
	return &FlightAPIClient{}
}

func (c *FlightAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	return nil, nil
}
