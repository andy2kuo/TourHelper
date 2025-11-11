package server

import (
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler 健康檢查處理器
type HealthCheckHandler struct {
	serviceName string
	env         string
	version     string
}

// NewHealthCheckHandler 建立健康檢查處理器
func NewHealthCheckHandler(serviceName, env, version string) *HealthCheckHandler {
	return &HealthCheckHandler{
		serviceName: serviceName,
		env:         env,
		version:     version,
	}
}

// Handle 處理健康檢查請求
func (h *HealthCheckHandler) Handle(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": h.serviceName,
		"env":     h.env,
		"version": h.version,
	})
}
