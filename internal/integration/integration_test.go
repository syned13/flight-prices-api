package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/syned13/flight-prices-api/internal/models"
)

type IntegrationTestSuite struct {
	suite.Suite
	baseURL string
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.baseURL = "http://localhost:8080"
	waitForHealthy(s.baseURL + "/health")
}

func waitForHealthy(url string) {
	timeout := time.After(60 * time.Second)
	tick := time.Tick(2 * time.Second)
	for {
		select {
		case <-timeout:
			panic("Service did not become healthy in time")
		case <-tick:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode < 500 {
				return
			}
		}
	}
}

func (s *IntegrationTestSuite) TestRegisterAndLogin() {
	client := &http.Client{}
	// Register user
	registerBody := `{"username":"integrationuser","password":"integrationpass"}`
	resp, err := client.Post(s.baseURL+"/register", "application/json", strings.NewReader(registerBody))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusCreated, resp.StatusCode)

	// Login user
	loginBody := `{"username":"integrationuser","password":"integrationpass"}`
	resp, err = client.Post(s.baseURL+"/login", "application/json", strings.NewReader(loginBody))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
}

func (s *IntegrationTestSuite) TestFetchItineraries() {
	client := &http.Client{}

	requestBody := `{"origin":"SDQ","destination":"JFK","departure_date":"2025-05-05T00:00:00Z","currency_code":"USD"}`
	resp, err := client.Post(s.baseURL+"/itineraries", "application/json", strings.NewReader(requestBody))
	s.Require().NoError(err)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()
	var result models.FlightSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	s.Require().Greater(len(result.Itineraries), 0, "should return at least one itinerary")
	s.Require().NotEmpty(result.Cheapest.Price.Total)
	s.Require().NotEmpty(result.Fastest.Price.Total)
}

func TestIntegrationTestSuite(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("skipping integration test; set INTEGRATION_TEST=1 to run")
	}
	suite.Run(t, new(IntegrationTestSuite))
}
