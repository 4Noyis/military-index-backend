package handler

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/4Noyis/military-index-backend/shared/auth"
	"github.com/4Noyis/military-index-backend/shared/utils"
)

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthHandler handles authentication endpoints
type AuthHandler struct{}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendBadRequest(w, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		utils.SendBadRequest(w, "Username and password are required")
		return
	}

	// Get credentials from environment
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminUsername == "" || adminPassword == "" {
		utils.SendInternalError(w, "Server configuration error")
		return
	}

	// Validate credentials
	if req.Username != adminUsername || req.Password != adminPassword {
		utils.SendUnauthorized(w, "Invalid username or password")
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(req.Username)
	if err != nil {
		utils.SendInternalError(w, "Failed to generate token")
		return
	}

	// Send response
	response := LoginResponse{Token: token}
	utils.SendSuccess(w, http.StatusOK, response, "Login successful")
}
