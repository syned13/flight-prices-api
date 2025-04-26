package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/controllers"
	"github.com/syned13/flight-prices-api/internal/repository/auth"
	flight_prices "github.com/syned13/flight-prices-api/internal/repository/itinerary-cache"
	auth_service "github.com/syned13/flight-prices-api/internal/services/auth"
	itinerary_fetcher "github.com/syned13/flight-prices-api/internal/services/itinerary-fetcher"
	"github.com/syned13/flight-prices-api/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	appConfig := config.GetConfig()
	if err := appConfig.Validate(); err != nil {
		log.Fatalf("failed to validate app config: %v", err)
	}

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.Mongo().URI()))
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	// Initialize Redis client
	redisAddr := appConfig.Redis().URI()
	if redisAddr == "" {
		log.Fatalf("failed to get redis address")
	}

	redisPassword := appConfig.Redis().Password()
	if redisPassword == "" {
		log.Fatalf("failed to get redis password")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	cacheRepo := flight_prices.NewItineraryCache(redisClient)

	// Initialize repositories
	authRepository := auth.NewAuthRepository(mongoClient)

	// Initialize services
	authService := auth_service.NewAuthService(authRepository)
	itineraryFetcherService := itinerary_fetcher.NewItineraryFetcherService(cacheRepo)

	router := mux.NewRouter()

	// Health endpoint for integration tests
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}).Methods("GET")

	// Initialize controllers
	authController := controllers.NewAuthController(router, authService)
	flightSearchController := controllers.NewFlightSearchController(router, itineraryFetcherService)

	// Register routes
	authController.RegisterRoutes()
	flightSearchController.RegisterRoutes()

	log.Printf("Starting server on port %s", appConfig.HTTPPort())
	err = http.ListenAndServe(fmt.Sprintf(":%s", appConfig.HTTPPort()), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
