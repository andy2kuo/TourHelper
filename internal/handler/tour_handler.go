package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/model"
	"github.com/andy2kuo/TourHelper/internal/service"
)

// TourHandler handles tour-related HTTP requests
type TourHandler struct {
	service service.TourService
	logger  *zap.Logger
}

// NewTourHandler creates a new tour handler
func NewTourHandler(service service.TourService, logger *zap.Logger) *TourHandler {
	return &TourHandler{
		service: service,
		logger:  logger,
	}
}

// CreateTour handles POST /api/v1/tours
func (h *TourHandler) CreateTour(c *gin.Context) {
	var req model.CreateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	tour, err := h.service.CreateTour(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create tour", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Success: true,
		Message: "Tour created successfully",
		Data:    tour,
	})
}

// GetTour handles GET /api/v1/tours/:id
func (h *TourHandler) GetTour(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "Invalid tour ID",
		})
		return
	}

	tour, err := h.service.GetTour(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get tour", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "Tour not found",
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    tour,
	})
}

// ListTours handles GET /api/v1/tours
func (h *TourHandler) ListTours(c *gin.Context) {
	filter := &model.TourFilter{
		Category: c.Query("category"),
		Country:  c.Query("country"),
		Budget:   c.Query("budget"),
		Season:   c.Query("season"),
	}

	// Parse limit and offset
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil {
			filter.Limit = l
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if o, err := strconv.Atoi(offset); err == nil {
			filter.Offset = o
		}
	}

	// Parse min_rating
	if minRating := c.Query("min_rating"); minRating != "" {
		if mr, err := strconv.ParseFloat(minRating, 64); err == nil {
			filter.MinRating = mr
		}
	}

	tours, total, err := h.service.ListTours(c.Request.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to list tours", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Calculate page number
	page := 1
	perPage := filter.Limit
	if perPage > 0 && filter.Offset > 0 {
		page = (filter.Offset / perPage) + 1
	}

	c.JSON(http.StatusOK, model.PaginatedResponse{
		Success: true,
		Data:    tours,
		Total:   total,
		Page:    page,
		PerPage: perPage,
	})
}

// UpdateTour handles PUT /api/v1/tours/:id
func (h *TourHandler) UpdateTour(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "Invalid tour ID",
		})
		return
	}

	var req model.UpdateTourRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	tour, err := h.service.UpdateTour(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("Failed to update tour", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Tour updated successfully",
		Data:    tour,
	})
}

// DeleteTour handles DELETE /api/v1/tours/:id
func (h *TourHandler) DeleteTour(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "Invalid tour ID",
		})
		return
	}

	if err := h.service.DeleteTour(c.Request.Context(), id); err != nil {
		h.logger.Error("Failed to delete tour", zap.Error(err), zap.Int64("id", id))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Tour deleted successfully",
	})
}

// SuggestTours handles POST /api/v1/tours/suggest
func (h *TourHandler) SuggestTours(c *gin.Context) {
	var preferences map[string]interface{}
	if err := c.ShouldBindJSON(&preferences); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	tours, err := h.service.SuggestTours(c.Request.Context(), preferences)
	if err != nil {
		h.logger.Error("Failed to suggest tours", zap.Error(err))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Tour suggestions generated",
		Data:    tours,
	})
}
