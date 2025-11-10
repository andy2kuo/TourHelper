package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/database"
	"github.com/andy2kuo/TourHelper/internal/handler"
	"github.com/andy2kuo/TourHelper/internal/middleware"
	"github.com/andy2kuo/TourHelper/internal/repository"
	"github.com/andy2kuo/TourHelper/internal/service"
	"github.com/andy2kuo/TourHelper/pkg/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := utils.NewLogger(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting TourHelper API server")

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database (optional - will fail gracefully if not available)
	var db *database.Database
	db, err = database.New(&cfg.Database, logger)
	if err != nil {
		logger.Warn("Failed to connect to database (will continue without DB)", zap.Error(err))
		// Continue without database - health checks will show degraded status
	}

	// Initialize repositories
	var tourRepo repository.TourRepository
	if db != nil {
		tourRepo = repository.NewTourRepository(db.DB, logger)
	}

	// Initialize services
	var tourService service.TourService
	if tourRepo != nil {
		tourService = service.NewTourService(tourRepo, logger)
	}

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(db, logger)
	var tourHandler *handler.TourHandler
	if tourService != nil {
		tourHandler = handler.NewTourHandler(tourService, logger)
	}

	// Setup router
	router := setupRouter(cfg, logger, healthHandler, tourHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting", zap.String("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close database connection
	if db != nil {
		if err := db.Close(); err != nil {
			logger.Error("Error closing database", zap.Error(err))
		}
	}

	logger.Info("Server stopped")
}

func setupRouter(cfg *config.Config, logger *zap.Logger, healthHandler *handler.HealthHandler, tourHandler *handler.TourHandler) *gin.Engine {
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.Logger(logger))
	router.Use(middleware.CORS())

	// Health check endpoints
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Tour endpoints
		if tourHandler != nil {
			tours := v1.Group("/tours")
			{
				tours.POST("", tourHandler.CreateTour)
				tours.GET("", tourHandler.ListTours)
				tours.GET("/:id", tourHandler.GetTour)
				tours.PUT("/:id", tourHandler.UpdateTour)
				tours.DELETE("/:id", tourHandler.DeleteTour)
				tours.POST("/suggest", tourHandler.SuggestTours)
			}
		}
	}

	return router
}
