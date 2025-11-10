package repository

import (
	"context"

	"github.com/andy2kuo/TourHelper/internal/model"
)

// TourRepository defines the interface for tour data access
type TourRepository interface {
	Create(ctx context.Context, tour *model.Tour) error
	GetByID(ctx context.Context, id int64) (*model.Tour, error)
	List(ctx context.Context, filter *model.TourFilter) ([]*model.Tour, error)
	Update(ctx context.Context, id int64, tour *model.Tour) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context, filter *model.TourFilter) (int64, error)
}
