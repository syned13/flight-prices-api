package controllers

import (
	"net/http"

	"github.com/gorilla/mux"
	itinerary_fetcher "github.com/syned13/flight-prices-api/internal/services/itinerary-fetcher"
)

type FlightSearchController struct {
	router                  *mux.Router
	itineraryFetcherService itinerary_fetcher.ItineraryFetcherService
}

func NewFlightSearchController(router *mux.Router, itineraryFetcherService itinerary_fetcher.ItineraryFetcherService) *FlightSearchController {
	return &FlightSearchController{
		router:                  router,
		itineraryFetcherService: itineraryFetcherService,
	}
}

func (c *FlightSearchController) RegisterRoutes() {
	c.router.HandleFunc("/flight-search", c.SearchFlights).Methods("GET")
}

func (c *FlightSearchController) SearchFlights(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
