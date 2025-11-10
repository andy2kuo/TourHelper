package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/database"
	"github.com/andy2kuo/TourHelper/internal/model"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db     *database.Database
	logger *zap.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// Health checks the health of the service
func (h *HealthHandler) Health(c *gin.Context) {
	dbStatus := "not_configured"
	status := "healthy"
	statusCode := http.StatusOK

	if h.db != nil {
		if err := h.db.Health(); err != nil {
			h.logger.Error("Database health check failed", zap.Error(err))
			dbStatus = "unhealthy"
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
		} else {
			dbStatus = "healthy"
		}
	}

	c.JSON(statusCode, model.HealthResponse{
		Status:   status,
		Database: dbStatus,
		Version:  "1.0.0",
	})
}

// Ready checks if the service is ready to serve requests
func (h *HealthHandler) Ready(c *gin.Context) {
	if h.db != nil {
		if err := h.db.Health(); err != nil {
			h.logger.Error("Database not ready", zap.Error(err))
			c.JSON(http.StatusServiceUnavailable, model.APIResponse{
				Success: false,
				Error:   "Service not ready - database unavailable",
			})
			return
		}
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Service is ready",
	})
}
