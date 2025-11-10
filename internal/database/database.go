package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/config"
)

// Database wraps the database connection
type Database struct {
	DB     *sql.DB
	logger *zap.Logger
}

// New creates a new database instance
func New(cfg *config.DatabaseConfig, logger *zap.Logger) (*Database, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	logger.Info("Database connection established",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.DBName),
	)

	return &Database{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.DB != nil {
		d.logger.Info("Closing database connection")
		return d.DB.Close()
	}
	return nil
}

// Health checks the database health
func (d *Database) Health() error {
	return d.DB.Ping()
}
