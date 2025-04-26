package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syned13/flight-prices-api/internal/middleware"
	"github.com/syned13/flight-prices-api/internal/models"
	"github.com/syned13/flight-prices-api/internal/services/auth"
)

type AuthController struct {
	router      *mux.Router
	authService auth.AuthService
}

func NewAuthController(router *mux.Router, authService auth.AuthService) *AuthController {
	return &AuthController{
		router:      router,
		authService: authService,
	}
}

func (c *AuthController) RegisterRoutes() {
	log.Printf("Registering auth routes")
	c.router.HandleFunc("/register", c.Register).Methods("POST")
	c.router.HandleFunc("/login", c.Login).Methods("POST")
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received login request")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := c.authService.Login(user)
	if err != nil {
		log.Printf("Login failed for user %s: %v", user.Username, err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	log.Printf("Login successful for user: %s", user.Username)
	middleware.WriteJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
	})
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received register request")
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := c.authService.Register(user); err != nil {
		log.Printf("Registration failed for user %s: %v", user.Username, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Registration successful for user: %s", user.Username)
	w.WriteHeader(http.StatusCreated)
}
