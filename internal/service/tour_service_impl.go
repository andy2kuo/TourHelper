package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/model"
	"github.com/andy2kuo/TourHelper/internal/repository"
)

// tourService implements TourService interface
type tourService struct {
	repo   repository.TourRepository
	logger *zap.Logger
}

// NewTourService creates a new tour service
func NewTourService(repo repository.TourRepository, logger *zap.Logger) TourService {
	return &tourService{
		repo:   repo,
		logger: logger,
	}
}

// CreateTour creates a new tour
func (s *tourService) CreateTour(ctx context.Context, req *model.CreateTourRequest) (*model.Tour, error) {
	// Validate business rules
	if err := s.validateCreateRequest(req); err != nil {
		return nil, err
	}

	tour := &model.Tour{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		Country:     req.Country,
		Category:    req.Category,
		Duration:    req.Duration,
		Season:      req.Season,
		Budget:      req.Budget,
		ImageURL:    req.ImageURL,
		Rating:      req.Rating,
	}

	if err := s.repo.Create(ctx, tour); err != nil {
		s.logger.Error("Failed to create tour in service", zap.Error(err))
		return nil, fmt.Errorf("failed to create tour: %w", err)
	}

	s.logger.Info("Tour created successfully", zap.Int64("id", tour.ID), zap.String("name", tour.Name))
	return tour, nil
}

// GetTour retrieves a tour by ID
func (s *tourService) GetTour(ctx context.Context, id int64) (*model.Tour, error) {
	tour, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get tour in service", zap.Error(err), zap.Int64("id", id))
		return nil, err
	}

	return tour, nil
}

// ListTours retrieves tours based on filter
func (s *tourService) ListTours(ctx context.Context, filter *model.TourFilter) ([]*model.Tour, int64, error) {
	// Set default pagination if not provided
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	tours, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list tours in service", zap.Error(err))
		return nil, 0, err
	}

	count, err := s.repo.Count(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to count tours in service", zap.Error(err))
		return tours, 0, nil // Return tours even if count fails
	}

	return tours, count, nil
}

// UpdateTour updates a tour
func (s *tourService) UpdateTour(ctx context.Context, id int64, req *model.UpdateTourRequest) (*model.Tour, error) {
	// Get existing tour
	existingTour, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if req.Name != nil {
		existingTour.Name = *req.Name
	}
	if req.Description != nil {
		existingTour.Description = *req.Description
	}
	if req.Location != nil {
		existingTour.Location = *req.Location
	}
	if req.Country != nil {
		existingTour.Country = *req.Country
	}
	if req.Category != nil {
		existingTour.Category = *req.Category
	}
	if req.Duration != nil {
		existingTour.Duration = *req.Duration
	}
	if req.Season != nil {
		existingTour.Season = *req.Season
	}
	if req.Budget != nil {
		existingTour.Budget = *req.Budget
	}
	if req.ImageURL != nil {
		existingTour.ImageURL = *req.ImageURL
	}
	if req.Rating != nil {
		existingTour.Rating = *req.Rating
	}

	if err := s.repo.Update(ctx, id, existingTour); err != nil {
		s.logger.Error("Failed to update tour in service", zap.Error(err), zap.Int64("id", id))
		return nil, fmt.Errorf("failed to update tour: %w", err)
	}

	s.logger.Info("Tour updated successfully", zap.Int64("id", id))
	return existingTour, nil
}

// DeleteTour deletes a tour
func (s *tourService) DeleteTour(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete tour in service", zap.Error(err), zap.Int64("id", id))
		return err
	}

	s.logger.Info("Tour deleted successfully", zap.Int64("id", id))
	return nil
}

// SuggestTours suggests tours based on user preferences
func (s *tourService) SuggestTours(ctx context.Context, preferences map[string]interface{}) ([]*model.Tour, error) {
	// Build filter from preferences
	filter := &model.TourFilter{
		Limit: 10, // Default to top 10 suggestions
	}

	if category, ok := preferences["category"].(string); ok {
		filter.Category = category
	}
	if country, ok := preferences["country"].(string); ok {
		filter.Country = country
	}
	if budget, ok := preferences["budget"].(string); ok {
		filter.Budget = budget
	}
	if season, ok := preferences["season"].(string); ok {
		filter.Season = season
	}
	if minRating, ok := preferences["min_rating"].(float64); ok {
		filter.MinRating = minRating
	}

	tours, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to suggest tours", zap.Error(err))
		return nil, fmt.Errorf("failed to suggest tours: %w", err)
	}

	s.logger.Info("Tours suggested successfully", zap.Int("count", len(tours)))
	return tours, nil
}

// validateCreateRequest validates the create tour request
func (s *tourService) validateCreateRequest(req *model.CreateTourRequest) error {
	// Add business validation rules here
	validBudgets := map[string]bool{"low": true, "medium": true, "high": true}
	if req.Budget != "" && !validBudgets[req.Budget] {
		return fmt.Errorf("invalid budget: must be 'low', 'medium', or 'high'")
	}

	if req.Rating < 0 || req.Rating > 5 {
		return fmt.Errorf("invalid rating: must be between 0 and 5")
	}

	return nil
}
