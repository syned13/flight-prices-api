package models

import "time"

type FlightSearchRequest struct {
	Origin        string
	Destination   string
	DepartureDate time.Time
}
