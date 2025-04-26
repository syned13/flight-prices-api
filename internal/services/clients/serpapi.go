package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type SerpAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type serpAPIResponse struct {
	BestFlights []struct {
		Flights []struct {
			DepartureAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"departure_airport"`
			ArrivalAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"arrival_airport"`
			Duration     int    `json:"duration"`
			Airplane     string `json:"airplane"`
			Airline      string `json:"airline"`
			FlightNumber string `json:"flight_number"`
		} `json:"flights"`
		Layovers []struct {
			Duration int    `json:"duration"`
			Name     string `json:"name"`
			ID       string `json:"id"`
		} `json:"layovers"`
		TotalDuration int     `json:"total_duration"`
		Price         float64 `json:"price"`
		Type          string  `json:"type"`
	} `json:"best_flights"`
	OtherFlights []struct {
		Flights []struct {
			DepartureAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"departure_airport"`
			ArrivalAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"arrival_airport"`
			Duration     int    `json:"duration"`
			Airplane     string `json:"airplane"`
			Airline      string `json:"airline"`
			FlightNumber string `json:"flight_number"`
		} `json:"flights"`
		Layovers []struct {
			Duration int    `json:"duration"`
			Name     string `json:"name"`
			ID       string `json:"id"`
		} `json:"layovers"`
		TotalDuration int     `json:"total_duration"`
		Price         float64 `json:"price"`
		Type          string  `json:"type"`
	} `json:"other_flights"`
}

func NewSerpAPIClient() *SerpAPIClient {
	return &SerpAPIClient{
		baseURL:    config.GetConfig().SerpAPI().BaseURL(),
		apiKey:     config.GetConfig().SerpAPI().APIKey(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *SerpAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("engine", "google_flights")
	params.Add("departure_id", request.Origin)
	params.Add("arrival_id", request.Destination)
	params.Add("outbound_date", request.DepartureDate.Format("2006-01-02"))
	params.Add("currency", request.CurrencyCode)
	params.Add("api_key", c.apiKey)

	// Create request
	reqURL := fmt.Sprintf("%s/search?%s", c.baseURL, params.Encode())
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var serpResp serpAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&serpResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var itineraries []models.Itinerary

	// Process best flights
	for _, flight := range serpResp.BestFlights {
		itinerary := models.Itinerary{
			Price: models.Price{
				Currency: request.CurrencyCode,
				Total:    strconv.FormatFloat(flight.Price, 'f', 2, 64),
			},
			DurationInMinutes: flight.TotalDuration,
			Stops:             len(flight.Layovers),
		}

		// Add segments
		for _, segment := range flight.Flights {
			departureTime, _ := time.Parse("2006-01-02 15:04", segment.DepartureAirport.Time)
			arrivalTime, _ := time.Parse("2006-01-02 15:04", segment.ArrivalAirport.Time)

			itinerary.Segments = append(itinerary.Segments, models.Segment{
				Airline:       segment.Airline,
				DepartureTime: departureTime,
				ArrivalTime:   arrivalTime,
				Carrier:       segment.Airline,
				Number:        segment.FlightNumber,
				Origin:        segment.DepartureAirport.ID,
				Destination:   segment.ArrivalAirport.ID,
			})
		}

		itineraries = append(itineraries, itinerary)
	}

	// Process other flights
	for _, flight := range serpResp.OtherFlights {
		itinerary := models.Itinerary{
			Price: models.Price{
				Currency: request.CurrencyCode,
				Total:    strconv.FormatFloat(flight.Price, 'f', 2, 64),
			},
			DurationInMinutes: flight.TotalDuration,
			Stops:             len(flight.Layovers),
		}

		// Add segments
		for _, segment := range flight.Flights {
			departureTime, _ := time.Parse("2006-01-02 15:04", segment.DepartureAirport.Time)
			arrivalTime, _ := time.Parse("2006-01-02 15:04", segment.ArrivalAirport.Time)

			itinerary.Segments = append(itinerary.Segments, models.Segment{
				Airline:       segment.Airline,
				DepartureTime: departureTime,
				ArrivalTime:   arrivalTime,
				Carrier:       segment.Airline,
				Number:        segment.FlightNumber,
				Origin:        segment.DepartureAirport.ID,
				Destination:   segment.ArrivalAirport.ID,
			})
		}

		itineraries = append(itineraries, itinerary)
	}

	return itineraries, nil
}
