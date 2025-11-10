package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/andy2kuo/TourHelper/internal/model"
)

// Recovery creates a middleware for recovering from panics
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Success: false,
					Error:   "Internal server error",
				})
				c.Abort()
			}
		}()

		c.Next()
	}
}
