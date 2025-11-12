package services

// WeatherService 天氣服務介面
type WeatherService interface {
	// TODO: 實作天氣相關業務邏輯
}

// weatherService 天氣服務實作
type weatherService struct {
	// 可能需要天氣 API 相關設定
}

// NewWeatherService 建立天氣服務
func NewWeatherService() WeatherService {
	return &weatherService{}
}
