package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/middleware"
	"github.com/syned13/flight-prices-api/internal/models"
)

type AuthController struct {
	router *mux.Router
}

func NewAuthController(router *mux.Router) *AuthController {
	return &AuthController{
		router: router,
	}
}

func (c *AuthController) RegisterRoutes() {
	log.Printf("Registering auth routes")
	c.router.HandleFunc("/login", c.Login).Methods("POST")
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received login request")
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidCredentials(loginReq.Username, loginReq.Password) {
		log.Printf("Invalid credentials for user: %s", loginReq.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(loginReq.Username)
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	log.Printf("Login successful for user: %s", loginReq.Username)
	middleware.WriteJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
	})
}

func isValidCredentials(username, password string) bool {
	return username == "admin" && password == "password"
}
