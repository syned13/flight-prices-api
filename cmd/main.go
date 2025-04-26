package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/controllers"
	itinerary_fetcher "github.com/syned13/flight-prices-api/internal/services/itinerary-fetcher"
	"github.com/syned13/flight-prices-api/pkg/config"
)

func getAppConfig() (*config.AppConfig, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		return nil, errors.New("HTTP_PORT is not set")
	}

	return &config.AppConfig{
		HttpPort: port,
	}, nil
}

func main() {
	appConfig, err := getAppConfig()
	if err != nil {
		log.Fatalf("failed to get app config: %v", err)
	}

	router := mux.NewRouter()

	itineraryFetcherService := itinerary_fetcher.NewItineraryFetcherService()
	flightSearchController := controllers.NewFlightSearchController(router, itineraryFetcherService)
	flightSearchController.RegisterRoutes()

	log.Printf("Starting server on port %s", appConfig.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", appConfig.HttpPort), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
