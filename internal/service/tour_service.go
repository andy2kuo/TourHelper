package service

import (
	"context"

	"github.com/andy2kuo/TourHelper/internal/model"
)

// TourService defines the interface for tour business logic
type TourService interface {
	CreateTour(ctx context.Context, req *model.CreateTourRequest) (*model.Tour, error)
	GetTour(ctx context.Context, id int64) (*model.Tour, error)
	ListTours(ctx context.Context, filter *model.TourFilter) ([]*model.Tour, int64, error)
	UpdateTour(ctx context.Context, id int64, req *model.UpdateTourRequest) (*model.Tour, error)
	DeleteTour(ctx context.Context, id int64) error
	SuggestTours(ctx context.Context, preferences map[string]interface{}) ([]*model.Tour, error)
}
