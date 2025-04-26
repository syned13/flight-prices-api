package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	DefaultHTTPPort         = "8080"
	DefaultAmadeusBaseURL   = "https://api.amadeus.com"
	DefaultFlightAPIBaseURL = "https://flightapi.com"
	DefaultSerpAPIBaseURL   = "https://serpapi.com"
	DefaultJWTExpiration    = "24h"
)

type AppConfig struct {
	httpPort  string
	amadeus   AmadeusConfig
	flightAPI FlightAPIConfig
	serpAPI   SerpAPIConfig
	jwt       JWTConfig
	mongo     MongoDBConfig
	redis     RedisConfig
}

type MongoDBConfig struct {
	uri      string
	database string
}

type RedisConfig struct {
	uri               string
	password          string
	cacheTTLInSeconds int
}

type AmadeusConfig struct {
	apiKey  string
	baseURL string
}

type FlightAPIConfig struct {
	apiKey  string
	baseURL string
}

type SerpAPIConfig struct {
	apiKey  string
	baseURL string
}

type JWTConfig struct {
	secret     string
	expiration string
}

var (
	instance *AppConfig
	once     sync.Once
)

func GetConfig() *AppConfig {
	once.Do(func() {
		instance = &AppConfig{
			httpPort: getEnvOrDefault("PORT", DefaultHTTPPort),
			amadeus: AmadeusConfig{
				apiKey:  getEnvOrDefault("AMADEUS_API_KEY", ""),
				baseURL: getEnvOrDefault("AMADEUS_BASE_URL", DefaultAmadeusBaseURL),
			},
			flightAPI: FlightAPIConfig{
				apiKey:  getEnvOrDefault("FLIGHT_API_KEY", ""),
				baseURL: getEnvOrDefault("FLIGHT_API_BASE_URL", DefaultFlightAPIBaseURL),
			},
			serpAPI: SerpAPIConfig{
				apiKey:  getEnvOrDefault("SERP_API_KEY", ""),
				baseURL: getEnvOrDefault("SERP_API_BASE_URL", DefaultSerpAPIBaseURL),
			},
			jwt: JWTConfig{
				secret:     os.Getenv("JWT_SECRET"),
				expiration: os.Getenv("JWT_EXPIRATION"),
			},
			mongo: MongoDBConfig{
				uri:      getEnvOrDefault("MONGO_URI", ""),
				database: getEnvOrDefault("MONGO_DATABASE", "flight-prices"),
			},
			redis: RedisConfig{
				uri:               getEnvOrDefault("REDIS_URI", ""),
				password:          getEnvOrDefault("REDIS_PASSWORD", ""),
				cacheTTLInSeconds: getIntEnvOrDefault("REDIS_CACHE_TTL_IN_SECONDS", 30),
			},
		}
	})
	return instance
}

func (c *AppConfig) HTTPPort() string { return c.httpPort }

func (c *AppConfig) Amadeus() AmadeusConfig     { return c.amadeus }
func (c *AppConfig) FlightAPI() FlightAPIConfig { return c.flightAPI }
func (c *AppConfig) SerpAPI() SerpAPIConfig     { return c.serpAPI }
func (c *AppConfig) JWT() JWTConfig             { return c.jwt }

func (c AmadeusConfig) APIKey() string  { return c.apiKey }
func (c AmadeusConfig) BaseURL() string { return c.baseURL }

func (c FlightAPIConfig) APIKey() string  { return c.apiKey }
func (c FlightAPIConfig) BaseURL() string { return c.baseURL }

func (c SerpAPIConfig) APIKey() string  { return c.apiKey }
func (c SerpAPIConfig) BaseURL() string { return c.baseURL }

func (c JWTConfig) Secret() string     { return c.secret }
func (c JWTConfig) Expiration() string { return c.expiration }

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		parsedValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}

		return parsedValue
	}

	return defaultValue
}

func (c *AppConfig) Validate() error {
	if c.HTTPPort() == "" {
		return fmt.Errorf("http port is not set")
	}
	if c.Amadeus().APIKey() == "" {
		return fmt.Errorf("amadeus API key is not set")
	}
	if c.FlightAPI().APIKey() == "" {
		return fmt.Errorf("flight API key is not set")
	}
	if c.SerpAPI().APIKey() == "" {
		return fmt.Errorf("serp API key is not set")
	}
	if c.Amadeus().BaseURL() == "" {
		return fmt.Errorf("amadeus base URL is not set")
	}
	if c.FlightAPI().BaseURL() == "" {
		return fmt.Errorf("flight API base URL is not set")
	}
	if c.SerpAPI().BaseURL() == "" {
		return fmt.Errorf("serp API base URL is not set")
	}

	if c.JWT().Secret() == "" {
		return fmt.Errorf("JWT secret is not set")
	}

	if c.JWT().Expiration() == "" {
		return fmt.Errorf("JWT expiration is not set")
	}

	return nil
}

func (c *AppConfig) Mongo() MongoDBConfig { return c.mongo }

func (c MongoDBConfig) URI() string      { return c.uri }
func (c MongoDBConfig) Database() string { return c.database }

func (c *AppConfig) Redis() RedisConfig { return c.redis }

func (c RedisConfig) URI() string            { return c.uri }
func (c RedisConfig) Password() string       { return c.password }
func (c RedisConfig) CacheTTLInSeconds() int { return c.cacheTTLInSeconds }
