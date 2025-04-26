package models

type FlightSearchResponse struct {
	Itineraries []Itinerary `json:"itineraries"`
	Cheapest    Itinerary   `json:"cheapest"`
	Fastest     Itinerary   `json:"fastest"`
}
