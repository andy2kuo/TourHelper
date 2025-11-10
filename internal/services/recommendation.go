package services

import (
	"math"

	"github.com/andy2kuo/TourHelper/internal/models"
)

// RecommendationService 推薦服務
type RecommendationService struct {
	weatherService *WeatherService
	// db *gorm.DB
}

// NewRecommendationService 建立推薦服務
func NewRecommendationService() *RecommendationService {
	return &RecommendationService{
		weatherService: NewWeatherService(),
	}
}

// GetRecommendations 取得旅遊推薦
func (s *RecommendationService) GetRecommendations(req models.RecommendationRequest) (*models.RecommendationResponse, error) {
	// 1. 取得當前天氣資訊
	weather, err := s.weatherService.GetWeather(req.Latitude, req.Longitude)
	if err != nil {
		// 天氣資訊取得失敗，使用預設值
		weather = &models.WeatherInfo{
			Temperature: 25.0,
			Condition:   "unknown",
			Humidity:    60,
			Description: "無法取得天氣資訊",
		}
	}

	// 2. 從資料庫取得附近的景點
	// TODO: 實作資料庫查詢
	destinations := s.getNearbyDestinations(req.Latitude, req.Longitude, req.MaxDistance)

	// 3. 根據天氣、距離、評分等因素計算適合度
	rankedDestinations := s.rankDestinations(destinations, req, weather)

	return &models.RecommendationResponse{
		Destinations: rankedDestinations,
		Weather:      *weather,
		Message:      s.generateRecommendationMessage(weather, len(rankedDestinations)),
	}, nil
}

// getNearbyDestinations 取得附近的景點
func (s *RecommendationService) getNearbyDestinations(lat, lon, maxDistance float64) []models.Destination {
	// TODO: 從資料庫查詢
	// 這裡先回傳範例資料
	return []models.Destination{
		{
			Name:        "陽明山國家公園",
			Description: "台北近郊的自然景觀勝地",
			Category:    "nature",
			Latitude:    25.1896,
			Longitude:   121.5453,
			Rating:      4.5,
			Address:     "台北市北投區竹子湖路",
			City:        "台北市",
			Region:      "北部",
			Country:     "Taiwan",
		},
		{
			Name:        "淡水老街",
			Description: "歷史悠久的老街，有許多特色小吃",
			Category:    "culture",
			Latitude:    25.1677,
			Longitude:   121.4406,
			Rating:      4.3,
			Address:     "新北市淡水區中正路",
			City:        "新北市",
			Region:      "北部",
			Country:     "Taiwan",
		},
	}
}

// rankDestinations 排序景點
func (s *RecommendationService) rankDestinations(
	destinations []models.Destination,
	req models.RecommendationRequest,
	weather *models.WeatherInfo,
) []models.DestinationWithDistance {
	var result []models.DestinationWithDistance

	for _, dest := range destinations {
		// 計算距離
		distance := calculateDistance(req.Latitude, req.Longitude, dest.Latitude, dest.Longitude)

		// 計算旅行時間（假設平均速度 50 km/h）
		travelTime := int(distance / 50 * 60)

		// 計算適合度評分
		suitability := s.calculateSuitability(dest, weather, distance)

		// 判斷天氣是否適合
		weatherMatch := s.isWeatherSuitable(dest.Category, weather.Condition)

		result = append(result, models.DestinationWithDistance{
			Destination:  dest,
			Distance:     distance,
			TravelTime:   travelTime,
			Suitability:  suitability,
			WeatherMatch: weatherMatch,
		})
	}

	// 根據適合度排序
	// TODO: 實作排序邏輯

	return result
}

// calculateDistance 計算兩點之間的距離（使用 Haversine 公式）
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // 地球半徑（公里）

	// 轉換為弧度
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// calculateSuitability 計算適合度評分
func (s *RecommendationService) calculateSuitability(
	dest models.Destination,
	weather *models.WeatherInfo,
	distance float64,
) float64 {
	score := 0.0

	// 評分因素：
	// 1. 景點評分 (40%)
	score += dest.Rating * 8

	// 2. 距離因素 (30%) - 距離越近分數越高
	if distance <= 10 {
		score += 30
	} else if distance <= 30 {
		score += 20
	} else if distance <= 50 {
		score += 10
	}

	// 3. 天氣適合度 (30%)
	if s.isWeatherSuitable(dest.Category, weather.Condition) {
		score += 30
	} else {
		score += 10
	}

	return score
}

// isWeatherSuitable 判斷天氣是否適合該類型景點
func (s *RecommendationService) isWeatherSuitable(category, condition string) bool {
	suitableMap := map[string][]string{
		"nature":    {"sunny", "cloudy"},
		"culture":   {"sunny", "cloudy", "rainy"},
		"food":      {"sunny", "cloudy", "rainy"},
		"shopping":  {"sunny", "cloudy", "rainy"},
		"adventure": {"sunny"},
	}

	if suitable, exists := suitableMap[category]; exists {
		for _, s := range suitable {
			if s == condition {
				return true
			}
		}
	}

	return false
}

// generateRecommendationMessage 產生推薦訊息
func (s *RecommendationService) generateRecommendationMessage(weather *models.WeatherInfo, count int) string {
	if count == 0 {
		return "很抱歉，目前找不到適合的景點推薦。"
	}

	weatherDesc := ""
	switch weather.Condition {
	case "sunny":
		weatherDesc = "今天天氣晴朗"
	case "cloudy":
		weatherDesc = "今天多雲"
	case "rainy":
		weatherDesc = "今天有雨"
	default:
		weatherDesc = "根據當前天氣"
	}

	return weatherDesc + "，為您推薦以下景點："
}
