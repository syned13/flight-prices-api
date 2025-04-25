package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/pkg/config"
)

func getAppConfig() (*config.AppConfig, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		return nil, errors.New("HTTP_PORT is not set")
	}

	return &config.AppConfig{
		HttpPort: port,
	}, nil
}

func main() {
	appConfig, err := getAppConfig()
	if err != nil {
		log.Fatalf("failed to get app config: %v", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	log.Printf("Starting server on port %s", appConfig.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", appConfig.HttpPort), router)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
