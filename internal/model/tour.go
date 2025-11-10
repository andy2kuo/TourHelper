package model

import "time"

// Tour represents a tour destination
type Tour struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Country     string    `json:"country"`
	Category    string    `json:"category"` // e.g., "beach", "mountain", "city", "cultural"
	Duration    int       `json:"duration"` // in days
	Season      string    `json:"season"`   // best season to visit
	Budget      string    `json:"budget"`   // "low", "medium", "high"
	ImageURL    string    `json:"image_url"`
	Rating      float64   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTourRequest represents the request to create a new tour
type CreateTourRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Location    string  `json:"location" binding:"required"`
	Country     string  `json:"country" binding:"required"`
	Category    string  `json:"category" binding:"required"`
	Duration    int     `json:"duration" binding:"required,min=1"`
	Season      string  `json:"season"`
	Budget      string  `json:"budget"`
	ImageURL    string  `json:"image_url"`
	Rating      float64 `json:"rating" binding:"min=0,max=5"`
}

// UpdateTourRequest represents the request to update a tour
type UpdateTourRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Location    *string  `json:"location"`
	Country     *string  `json:"country"`
	Category    *string  `json:"category"`
	Duration    *int     `json:"duration"`
	Season      *string  `json:"season"`
	Budget      *string  `json:"budget"`
	ImageURL    *string  `json:"image_url"`
	Rating      *float64 `json:"rating"`
}

// TourFilter represents filter criteria for listing tours
type TourFilter struct {
	Category string
	Country  string
	Budget   string
	Season   string
	MinRating float64
	Limit    int
	Offset   int
}
