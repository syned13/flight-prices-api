package flight_prices

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type ItineraryCache interface {
	GetItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error)
	SaveItineraries(ctx context.Context, request models.FlightSearchRequest, itineraries []models.Itinerary) error
}

type flightPricesRepository struct {
	redis *redis.Client
}

func NewItineraryCache(redis *redis.Client) ItineraryCache {
	return &flightPricesRepository{
		redis: redis,
	}
}

func (r *flightPricesRepository) GetItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	key := fmt.Sprintf("itineraries:%s:%s:%s", request.Origin, request.Destination, request.DepartureDate.Format("2006-01-02"))

	itinerariesJSON, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var itineraries []models.Itinerary
	err = json.Unmarshal([]byte(itinerariesJSON), &itineraries)
	if err != nil {
		return nil, err
	}

	return itineraries, nil
}

func (r *flightPricesRepository) SaveItineraries(ctx context.Context, request models.FlightSearchRequest, itineraries []models.Itinerary) error {
	key := fmt.Sprintf("itineraries:%s:%s:%s", request.Origin, request.Destination, request.DepartureDate.Format("2006-01-02"))

	itinerariesJSON, err := json.Marshal(itineraries)
	if err != nil {
		return err
	}

	expiration := time.Duration(config.GetConfig().Redis().CacheTTLInSeconds()) * time.Second

	return r.redis.Set(ctx, key, itinerariesJSON, expiration).Err()
}
