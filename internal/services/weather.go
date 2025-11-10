package services

import (
	"fmt"

	"github.com/andy2kuo/TourHelper/internal/models"
)

// WeatherService 天氣服務
type WeatherService struct {
	apiKey   string
	provider string
}

// NewWeatherService 建立天氣服務
func NewWeatherService() *WeatherService {
	// TODO: 從設定檔讀取
	return &WeatherService{
		apiKey:   "",
		provider: "openweathermap",
	}
}

// GetWeather 取得天氣資訊
func (s *WeatherService) GetWeather(lat, lon float64) (*models.WeatherInfo, error) {
	// TODO: 實作實際的天氣 API 呼叫
	// 這裡先回傳模擬資料

	if s.apiKey == "" {
		return nil, fmt.Errorf("天氣 API 金鑰未設定")
	}

	// 模擬天氣資料
	return &models.WeatherInfo{
		Temperature: 25.0,
		Condition:   "sunny",
		Humidity:    60,
		Description: "晴朗",
	}, nil
}

// GetWeatherForecast 取得天氣預報
func (s *WeatherService) GetWeatherForecast(lat, lon float64, days int) ([]models.WeatherInfo, error) {
	// TODO: 實作天氣預報功能
	return nil, fmt.Errorf("天氣預報功能尚未實作")
}
