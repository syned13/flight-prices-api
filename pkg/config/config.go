package config

type AppConfig struct {
	HttpPort  string
	Amadeus   AmadeusConfig
	FlightAPI FlightAPIConfig
	SerpAPI   SerpAPIConfig
}

type AmadeusConfig struct {
	APIKey  string
	BaseURL string
}

type FlightAPIConfig struct {
	APIKey  string
	BaseURL string
}

type SerpAPIConfig struct {
	APIKey  string
	BaseURL string
}
