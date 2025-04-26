package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/controllers"
	itinerary_fetcher "github.com/syned13/flight-prices-api/internal/services/itinerary-fetcher"
	"github.com/syned13/flight-prices-api/pkg/config"
)

func main() {
	appConfig := config.GetConfig()
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("failed to validate app config: %v", err)
	}

	router := mux.NewRouter()

	itineraryFetcherService := itinerary_fetcher.NewItineraryFetcherService()
	flightSearchController := controllers.NewFlightSearchController(router, itineraryFetcherService)
	flightSearchController.RegisterRoutes()

	log.Printf("Starting server on port %s", appConfig.HTTPPort())
	err := http.ListenAndServe(fmt.Sprintf(":%s", appConfig.HTTPPort()), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
