package itinerary_fetcher

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syned13/flight-prices-api/internal/models"
	cachemock "github.com/syned13/flight-prices-api/mocks/repository/itinerary-cache"
	clientmock "github.com/syned13/flight-prices-api/mocks/services/clients"
	"go.uber.org/mock/gomock"
)

func TestItineraryFetcherService_FetchItineraries(t *testing.T) {
	t.Run("should fetch itineraries", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cache := cachemock.NewMockItineraryCache(ctrl)
		cache.EXPECT().GetItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)
		cache.EXPECT().SaveItineraries(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		amadeusClient := clientmock.NewMockItineraryFetcherClient(ctrl)
		flightAPIClient := clientmock.NewMockItineraryFetcherClient(ctrl)
		serpAPIClient := clientmock.NewMockItineraryFetcherClient(ctrl)

		expectedItineraries := []models.Itinerary{
			{
				Price:             models.Price{Currency: "USD", Total: "100.00"},
				DurationInMinutes: 120,
				Stops:             0,
			},
		}
		amadeusClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return(expectedItineraries, nil)
		flightAPIClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)
		serpAPIClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)

		service := &itineraryFetcherService{
			amadeusClient:   amadeusClient,
			flightAPIClient: flightAPIClient,
			serpAPIClient:   serpAPIClient,
			cache:           cache,
		}

		resp, err := service.FetchItineraries(context.Background(), models.FlightSearchRequest{})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, 1, len(resp.Itineraries))
		assert.Equal(t, "100.00", resp.Cheapest.Price.Total)
		assert.Equal(t, "100.00", resp.Fastest.Price.Total)
	})

	t.Run("should return error if no itineraries are found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cache := cachemock.NewMockItineraryCache(ctrl)
		cache.EXPECT().GetItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)

		amadeusClient := clientmock.NewMockItineraryFetcherClient(ctrl)
		flightAPIClient := clientmock.NewMockItineraryFetcherClient(ctrl)
		serpAPIClient := clientmock.NewMockItineraryFetcherClient(ctrl)

		amadeusClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)
		flightAPIClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)
		serpAPIClient.EXPECT().FetchItineraries(gomock.Any(), gomock.Any()).Return([]models.Itinerary{}, nil)

		service := &itineraryFetcherService{
			amadeusClient:   amadeusClient,
			flightAPIClient: flightAPIClient,
			serpAPIClient:   serpAPIClient,
			cache:           cache,
		}

		resp, err := service.FetchItineraries(context.Background(), models.FlightSearchRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "no itineraries found", err.Error())
	})
}
