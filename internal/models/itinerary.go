package models

import "time"

type Itinerary struct {
	Price             Price     `json:"price"`
	DurationInMinutes int       `json:"durationInMinutes"`
	Segments          []Segment `json:"segments"`
	Stops             int       `json:"stops"`
}

type Segment struct {
	Airline       string    `json:"airline"`
	DepartureTime time.Time `json:"departureTime"`
	ArrivalTime   time.Time `json:"arrivalTime"`
	Carrier       string    `json:"carrier"`
	Number        string    `json:"number"`
	Origin        string    `json:"origin"`
	Destination   string    `json:"destination"`
}

type Price struct {
	Currency string `json:"currency"`
	Total    string `json:"total"`
}
