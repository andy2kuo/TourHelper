package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/andy2kuo/TourHelper/internal/models"
)

// GetRecommendations 取得旅遊推薦
func GetRecommendations(c *gin.Context) {
	var req models.RecommendationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求格式",
			"details": err.Error(),
		})
		return
	}

	// TODO: 實作推薦邏輯
	// recommendationService := services.NewRecommendationService()
	// recommendations, err := recommendationService.GetRecommendations(req)

	// 暫時回傳範例資料
	c.JSON(http.StatusOK, models.RecommendationResponse{
		Destinations: []models.DestinationWithDistance{},
		Weather: models.WeatherInfo{
			Temperature: 25.0,
			Condition:   "sunny",
			Humidity:    60,
			Description: "晴朗",
		},
		Message: "推薦功能開發中",
	})
}

// UpdatePreferences 更新使用者偏好
func UpdatePreferences(c *gin.Context) {
	var prefs models.UserPreferences

	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "無效的請求格式",
			"details": err.Error(),
		})
		return
	}

	// TODO: 實作偏好更新邏輯

	c.JSON(http.StatusOK, gin.H{
		"message": "偏好設定已更新",
		"preferences": prefs,
	})
}

// GetPreferences 取得使用者偏好
func GetPreferences(c *gin.Context) {
	// TODO: 從資料庫取得使用者偏好

	// 暫時回傳預設值
	c.JSON(http.StatusOK, models.UserPreferences{
		MaxDistance:       50.0,
		PreferredWeather:  "any",
		PreferredCategory: "nature",
		MinRating:         3.0,
		Budget:            "medium",
	})
}

// HandleWebSocket 處理 WebSocket 連線
func HandleWebSocket(c *gin.Context) {
	// TODO: 實作 WebSocket 處理邏輯
	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "WebSocket 功能開發中",
	})
}
