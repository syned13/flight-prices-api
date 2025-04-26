package controllers

import (
	"encoding/json"
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
	c.router.HandleFunc("/login", c.Login).Methods("POST")
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !isValidCredentials(loginReq.Username, loginReq.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := middleware.GenerateToken(loginReq.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
	})
}

// This is a simple mock function
func isValidCredentials(username, password string) bool {
	return username == "admin" && password == "password"
}
