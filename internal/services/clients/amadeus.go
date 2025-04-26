package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/pkg/config"
)

type AmadeusClient struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	accessToken string
}

type amadeusAuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type amadeusFlightSearchRequest struct {
	CurrencyCode       string `json:"currencyCode"`
	OriginDestinations []struct {
		ID            string `json:"id"`
		OriginCode    string `json:"originLocationCode"`
		DestCode      string `json:"destinationLocationCode"`
		DepartureDate string `json:"departureDateTimeRange"`
	} `json:"originDestinations"`
	Travelers []struct {
		ID   string `json:"id"`
		Type string `json:"travelerType"`
	} `json:"travelers"`
	Sources []string `json:"sources"`
}

type amadeusFlightSearchResponse struct {
	Data []struct {
		Price struct {
			Total    string `json:"total"`
			Currency string `json:"currency"`
		} `json:"price"`
		Itineraries []struct {
			Duration string `json:"duration"`
			Segments []struct {
				Departure struct {
					IATACode string `json:"iataCode"`
					At       string `json:"at"`
				} `json:"departure"`
				Arrival struct {
					IATACode string `json:"iataCode"`
					At       string `json:"at"`
				} `json:"arrival"`
				CarrierCode string `json:"carrierCode"`
				Number      string `json:"number"`
			} `json:"segments"`
		} `json:"itineraries"`
	} `json:"data"`
}

func NewAmadeusClient() *AmadeusClient {
	return &AmadeusClient{
		baseURL:    config.GetConfig().Amadeus().BaseURL(),
		apiKey:     config.GetConfig().Amadeus().APIKey(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *AmadeusClient) authenticate(ctx context.Context) error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.apiKey)
	data.Set("client_secret", c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/security/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth request failed with status: %d", resp.StatusCode)
	}

	var authResp amadeusAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.accessToken = authResp.AccessToken
	return nil
}

func (c *AmadeusClient) FetchItineraries(ctx context.Context, request models.FlightSearchRequest) ([]models.Itinerary, error) {
	if c.accessToken == "" {
		if err := c.authenticate(ctx); err != nil {
			return nil, fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	searchReq := amadeusFlightSearchRequest{
		CurrencyCode: request.CurrencyCode,
		OriginDestinations: []struct {
			ID            string `json:"id"`
			OriginCode    string `json:"originLocationCode"`
			DestCode      string `json:"destinationLocationCode"`
			DepartureDate string `json:"departureDateTimeRange"`
		}{
			{
				ID:            "1",
				OriginCode:    request.Origin,
				DestCode:      request.Destination,
				DepartureDate: request.DepartureDate.Format("2006-01-02"),
			},
		},
		Travelers: []struct {
			ID   string `json:"id"`
			Type string `json:"travelerType"`
		}{
			{
				ID:   "1",
				Type: "ADULT",
			},
		},
		Sources: []string{"GDS"},
	}

	jsonData, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/shopping/flight-offers", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+c.accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status: %d", resp.StatusCode)
	}

	var searchResp amadeusFlightSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	var itineraries []models.Itinerary
	for _, data := range searchResp.Data {
		for _, itin := range data.Itineraries {
			itinerary := models.Itinerary{
				Price: models.Price{
					Currency: data.Price.Currency,
					Total:    data.Price.Total,
				},
				Duration: itin.Duration,
				Stops:    len(itin.Segments) - 1,
			}

			for _, seg := range itin.Segments {
				departure, _ := time.Parse(time.RFC3339, seg.Departure.At)
				arrival, _ := time.Parse(time.RFC3339, seg.Arrival.At)

				segment := models.Segment{
					Airline:       seg.CarrierCode,
					DepartureTime: departure,
					ArrivalTime:   arrival,
					Carrier:       seg.CarrierCode,
					Number:        seg.Number,
					Origin:        seg.Departure.IATACode,
					Destination:   seg.Arrival.IATACode,
				}
				itinerary.Segments = append(itinerary.Segments, segment)
			}
			itineraries = append(itineraries, itinerary)
		}
	}

	return itineraries, nil
}
