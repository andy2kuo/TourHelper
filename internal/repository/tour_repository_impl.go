package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/model"
)

// tourRepository implements TourRepository interface
type tourRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewTourRepository creates a new tour repository
func NewTourRepository(db *sql.DB, logger *zap.Logger) TourRepository {
	return &tourRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new tour
func (r *tourRepository) Create(ctx context.Context, tour *model.Tour) error {
	query := `
		INSERT INTO tours (name, description, location, country, category, duration, season, budget, image_url, rating, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`

	now := time.Now()
	tour.CreatedAt = now
	tour.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		query,
		tour.Name,
		tour.Description,
		tour.Location,
		tour.Country,
		tour.Category,
		tour.Duration,
		tour.Season,
		tour.Budget,
		tour.ImageURL,
		tour.Rating,
		tour.CreatedAt,
		tour.UpdatedAt,
	).Scan(&tour.ID)

	if err != nil {
		r.logger.Error("Failed to create tour", zap.Error(err))
		return fmt.Errorf("failed to create tour: %w", err)
	}

	return nil
}

// GetByID retrieves a tour by ID
func (r *tourRepository) GetByID(ctx context.Context, id int64) (*model.Tour, error) {
	query := `
		SELECT id, name, description, location, country, category, duration, season, budget, image_url, rating, created_at, updated_at
		FROM tours
		WHERE id = $1
	`

	tour := &model.Tour{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tour.ID,
		&tour.Name,
		&tour.Description,
		&tour.Location,
		&tour.Country,
		&tour.Category,
		&tour.Duration,
		&tour.Season,
		&tour.Budget,
		&tour.ImageURL,
		&tour.Rating,
		&tour.CreatedAt,
		&tour.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("tour not found")
	}
	if err != nil {
		r.logger.Error("Failed to get tour", zap.Error(err), zap.Int64("id", id))
		return nil, fmt.Errorf("failed to get tour: %w", err)
	}

	return tour, nil
}

// List retrieves tours based on filter
func (r *tourRepository) List(ctx context.Context, filter *model.TourFilter) ([]*model.Tour, error) {
	query := `
		SELECT id, name, description, location, country, category, duration, season, budget, image_url, rating, created_at, updated_at
		FROM tours
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filter.Category)
		argCount++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argCount)
		args = append(args, filter.Country)
		argCount++
	}

	if filter.Budget != "" {
		query += fmt.Sprintf(" AND budget = $%d", argCount)
		args = append(args, filter.Budget)
		argCount++
	}

	if filter.Season != "" {
		query += fmt.Sprintf(" AND season = $%d", argCount)
		args = append(args, filter.Season)
		argCount++
	}

	if filter.MinRating > 0 {
		query += fmt.Sprintf(" AND rating >= $%d", argCount)
		args = append(args, filter.MinRating)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
		argCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to list tours", zap.Error(err))
		return nil, fmt.Errorf("failed to list tours: %w", err)
	}
	defer rows.Close()

	tours := []*model.Tour{}
	for rows.Next() {
		tour := &model.Tour{}
		err := rows.Scan(
			&tour.ID,
			&tour.Name,
			&tour.Description,
			&tour.Location,
			&tour.Country,
			&tour.Category,
			&tour.Duration,
			&tour.Season,
			&tour.Budget,
			&tour.ImageURL,
			&tour.Rating,
			&tour.CreatedAt,
			&tour.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan tour", zap.Error(err))
			return nil, fmt.Errorf("failed to scan tour: %w", err)
		}
		tours = append(tours, tour)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tours: %w", err)
	}

	return tours, nil
}

// Update updates a tour
func (r *tourRepository) Update(ctx context.Context, id int64, tour *model.Tour) error {
	query := `
		UPDATE tours
		SET name = $1, description = $2, location = $3, country = $4, category = $5, 
		    duration = $6, season = $7, budget = $8, image_url = $9, rating = $10, updated_at = $11
		WHERE id = $12
	`

	tour.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		tour.Name,
		tour.Description,
		tour.Location,
		tour.Country,
		tour.Category,
		tour.Duration,
		tour.Season,
		tour.Budget,
		tour.ImageURL,
		tour.Rating,
		tour.UpdatedAt,
		id,
	)

	if err != nil {
		r.logger.Error("Failed to update tour", zap.Error(err), zap.Int64("id", id))
		return fmt.Errorf("failed to update tour: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found")
	}

	return nil
}

// Delete deletes a tour
func (r *tourRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tours WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete tour", zap.Error(err), zap.Int64("id", id))
		return fmt.Errorf("failed to delete tour: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tour not found")
	}

	return nil
}

// Count counts tours based on filter
func (r *tourRepository) Count(ctx context.Context, filter *model.TourFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM tours WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argCount)
		args = append(args, filter.Category)
		argCount++
	}

	if filter.Country != "" {
		query += fmt.Sprintf(" AND country = $%d", argCount)
		args = append(args, filter.Country)
		argCount++
	}

	if filter.Budget != "" {
		query += fmt.Sprintf(" AND budget = $%d", argCount)
		args = append(args, filter.Budget)
		argCount++
	}

	if filter.Season != "" {
		query += fmt.Sprintf(" AND season = $%d", argCount)
		args = append(args, filter.Season)
		argCount++
	}

	if filter.MinRating > 0 {
		query += fmt.Sprintf(" AND rating >= $%d", argCount)
		args = append(args, filter.MinRating)
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to count tours", zap.Error(err))
		return 0, fmt.Errorf("failed to count tours: %w", err)
	}

	return count, nil
}

// Helper function to build WHERE clause (not used currently but useful for complex queries)
func buildWhereClause(filter *model.TourFilter) (string, []interface{}) {
	conditions := []string{}
	args := []interface{}{}
	argCount := 1

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, filter.Category)
		argCount++
	}

	if filter.Country != "" {
		conditions = append(conditions, fmt.Sprintf("country = $%d", argCount))
		args = append(args, filter.Country)
		argCount++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	return whereClause, args
}
