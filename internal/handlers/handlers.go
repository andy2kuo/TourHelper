package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/andy2kuo/TourHelper/config"
	"github.com/andy2kuo/TourHelper/internal/auth"
	"github.com/andy2kuo/TourHelper/internal/models"
	ws "github.com/andy2kuo/TourHelper/internal/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Handler manages HTTP handlers
type Handler struct {
	config    *config.Config
	authStore *auth.Store
	wsHub     *ws.Hub
}

// NewHandler creates a new handler
func NewHandler(cfg *config.Config, authStore *auth.Store, wsHub *ws.Hub) *Handler {
	return &Handler{
		config:    cfg,
		authStore: authStore,
		wsHub:     wsHub,
	}
}

// HandleLogin handles user login
func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate user
	user, err := h.authStore.Authenticate(req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user, h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Send response
	response := models.LoginResponse{
		Token:    token,
		Username: user.Username,
		Platform: user.Platform,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Get token from query parameter or header
	token := r.URL.Query().Get("token")
	if token == "" {
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
	}

	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// Validate token
	claims, err := auth.ValidateToken(token, h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Serve WebSocket
	ws.ServeWS(h.wsHub, conn, claims.Username, claims.Platform)
}

// HandleGetConfig returns configuration for the client
func (h *Handler) HandleGetConfig(w http.ResponseWriter, r *http.Request) {
	// Validate token
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	_, err := auth.ValidateToken(token, h.config.JWTSecret)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Return config
	config := map[string]interface{}{
		"googleMapsApiKey": h.config.GoogleMapsAPIKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

// HandleHealth handles health check
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

// EnableCORS enables CORS for all handlers
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
