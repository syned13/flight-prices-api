package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/middleware"
	"github.com/syned13/flight-prices-api/internal/models"
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
	log.Printf("Registering flight search routes")
	c.router.HandleFunc("/flights/search", middleware.AuthMiddleware(c.SearchFlights)).Methods("GET")
}

func (c *FlightSearchController) SearchFlights(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	dateStr := r.URL.Query().Get("date")

	if origin == "" || destination == "" || dateStr == "" {
		http.Error(w, "Missing required parameters: origin, destination, and date are required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	request := models.FlightSearchRequest{
		Origin:        origin,
		Destination:   destination,
		DepartureDate: date,
		CurrencyCode:  "USD",
	}

	response, err := c.itineraryFetcherService.FetchItineraries(r.Context(), request)
	if err != nil {
		log.Printf("Error fetching itineraries: %v", err)
		http.Error(w, "Failed to fetch flight prices", http.StatusInternalServerError)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, response)
}
