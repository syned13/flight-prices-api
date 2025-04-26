package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type FlightAPIClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type flightAPIResponse struct {
	Itineraries []struct {
		ID             string   `json:"id"`
		LegIDs         []string `json:"leg_ids"`
		PricingOptions []struct {
			Price struct {
				Amount       float64 `json:"amount"`
				UpdateStatus string  `json:"update_status"`
				LastUpdated  string  `json:"last_updated"`
			} `json:"price"`
		} `json:"pricing_options"`
	} `json:"itineraries"`
	Legs []struct {
		ID         string   `json:"id"`
		Departure  string   `json:"departure"`
		Arrival    string   `json:"arrival"`
		Duration   int      `json:"duration"`
		StopCount  int      `json:"stop_count"`
		SegmentIDs []string `json:"segment_ids"`
	} `json:"legs"`
	Segments []struct {
		ID                    string `json:"id"`
		Departure             string `json:"departure"`
		Arrival               string `json:"arrival"`
		Duration              int    `json:"duration"`
		MarketingFlightNumber string `json:"marketing_flight_number"`
		MarketingCarrierID    int    `json:"marketing_carrier_id"`
		OperatingCarrierID    int    `json:"operating_carrier_id"`
	} `json:"segments"`
}

func NewFlightAPIClient() *FlightAPIClient {
	return &FlightAPIClient{
		baseURL:    config.GetConfig().FlightAPI().BaseURL(),
		apiKey:     config.GetConfig().FlightAPI().APIKey(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *FlightAPIClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	url := fmt.Sprintf("%s/onewaytrip/%s/%s/%s/%s/1/0/0/Economy/%s",
		c.baseURL,
		c.apiKey,
		request.Origin,
		request.Destination,
		request.DepartureDate.Format("2006-01-02"),
		request.CurrencyCode,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	var flightResp flightAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&flightResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	legMap := make(map[string]struct {
		ID         string   `json:"id"`
		Departure  string   `json:"departure"`
		Arrival    string   `json:"arrival"`
		Duration   int      `json:"duration"`
		StopCount  int      `json:"stop_count"`
		SegmentIDs []string `json:"segment_ids"`
	})

	for i := range flightResp.Legs {
		legMap[flightResp.Legs[i].ID] = flightResp.Legs[i]
	}

	segmentMap := make(map[string]struct {
		ID                    string `json:"id"`
		Departure             string `json:"departure"`
		Arrival               string `json:"arrival"`
		Duration              int    `json:"duration"`
		MarketingFlightNumber string `json:"marketing_flight_number"`
		MarketingCarrierID    int    `json:"marketing_carrier_id"`
		OperatingCarrierID    int    `json:"operating_carrier_id"`
	})

	for i := range flightResp.Segments {
		segmentMap[flightResp.Segments[i].ID] = flightResp.Segments[i]
	}

	var itineraries []models.Itinerary
	for _, itin := range flightResp.Itineraries {
		if len(itin.PricingOptions) == 0 {
			continue
		}

		if len(itin.LegIDs) == 0 {
			continue
		}

		leg, ok := legMap[itin.LegIDs[0]]
		if !ok {
			continue
		}

		itinerary := models.Itinerary{
			Price: models.Price{
				Currency: request.CurrencyCode,
				Total:    fmt.Sprintf("%.2f", itin.PricingOptions[0].Price.Amount),
			},
			DurationInMinutes: leg.Duration,
			Stops:             leg.StopCount,
		}

		for _, segID := range leg.SegmentIDs {
			seg, ok := segmentMap[segID]
			if !ok {
				continue
			}

			departure, _ := time.Parse(time.RFC3339, seg.Departure)
			arrival, _ := time.Parse(time.RFC3339, seg.Arrival)

			segment := models.Segment{
				Airline:       fmt.Sprintf("%d", seg.MarketingCarrierID),
				DepartureTime: departure,
				ArrivalTime:   arrival,
				Carrier:       fmt.Sprintf("%d", seg.OperatingCarrierID),
				Number:        seg.MarketingFlightNumber,
				Origin:        request.Origin,
				Destination:   request.Destination,
			}
			itinerary.Segments = append(itinerary.Segments, segment)
		}

		itineraries = append(itineraries, itinerary)
	}

	return itineraries, nil
}
