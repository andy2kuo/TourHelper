package models

import "time"

// User represents a user in the system
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"` // Never expose password in JSON
	Email        string    `json:"email"`
	Platform     string    `json:"platform"` // web, line, telegram
	CreatedAt    time.Time `json:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Platform string `json:"platform"` // web, line, telegram
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Platform string `json:"platform"`
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
}

// Location represents a geographical location
type Location struct {
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
	Name string  `json:"name"`
}

// TourSuggestion represents a tour suggestion
type TourSuggestion struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Location    Location   `json:"location"`
	CreatedAt   time.Time  `json:"created_at"`
}
