package models

import (
	"time"

	"gorm.io/gorm"
)

// User 使用者模型
type User struct {
	gorm.Model
	ExternalID   string `gorm:"uniqueIndex;not null"` // Line ID 或 Telegram ID
	Platform     string `gorm:"not null"`             // line, telegram, web
	Username     string
	DisplayName  string
	Preferences  UserPreferences `gorm:"foreignKey:UserID"`
	SearchHistory []SearchHistory `gorm:"foreignKey:UserID"`
}

// UserPreferences 使用者偏好設定
type UserPreferences struct {
	gorm.Model
	UserID            uint
	MaxDistance       float64 `gorm:"default:50"`    // 最大距離（公里）
	PreferredWeather  string  `gorm:"default:any"`   // sunny, cloudy, rainy, any
	PreferredCategory string  // nature, culture, food, shopping, adventure
	MinRating         float64 `gorm:"default:3.0"`   // 最低評分
	Budget            string  `gorm:"default:medium"` // low, medium, high
}

// Destination 旅遊目的地
type Destination struct {
	gorm.Model
	Name        string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Category    string  // nature, culture, food, shopping, adventure
	Latitude    float64 `gorm:"not null"`
	Longitude   float64 `gorm:"not null"`
	Rating      float64 `gorm:"default:0"`
	Address     string
	City        string
	Region      string
	Country     string  `gorm:"default:Taiwan"`
	ImageURL    string
	Website     string
	Tags        []Tag   `gorm:"many2many:destination_tags;"`
}

// Tag 標籤
type Tag struct {
	gorm.Model
	Name         string        `gorm:"uniqueIndex;not null"`
	Destinations []Destination `gorm:"many2many:destination_tags;"`
}

// SearchHistory 搜尋歷史記錄
type SearchHistory struct {
	gorm.Model
	UserID           uint
	SearchLatitude   float64
	SearchLongitude  float64
	SearchLocation   string
	Weather          string
	RecommendationID uint
	Destination      Destination `gorm:"foreignKey:RecommendationID"`
	Clicked          bool        `gorm:"default:false"` // 是否點擊查看詳情
	Visited          bool        `gorm:"default:false"` // 是否標記為已造訪
}

// WeatherData 天氣資料快取
type WeatherData struct {
	gorm.Model
	Latitude    float64
	Longitude   float64
	Temperature float64
	Condition   string // sunny, cloudy, rainy, snowy
	Humidity    int
	WindSpeed   float64
	FetchedAt   time.Time
	ExpireAt    time.Time
}

// RecommendationRequest 推薦請求結構
type RecommendationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	UserID    uint    `json:"user_id,omitempty"`
	Category  string  `json:"category,omitempty"`
	MaxDistance float64 `json:"max_distance,omitempty"`
}

// RecommendationResponse 推薦回應結構
type RecommendationResponse struct {
	Destinations []DestinationWithDistance `json:"destinations"`
	Weather      WeatherInfo               `json:"weather"`
	Message      string                    `json:"message,omitempty"`
}

// DestinationWithDistance 帶距離資訊的目的地
type DestinationWithDistance struct {
	Destination
	Distance      float64 `json:"distance"`       // 距離（公里）
	TravelTime    int     `json:"travel_time"`    // 預估旅行時間（分鐘）
	Suitability   float64 `json:"suitability"`    // 適合度評分 (0-100)
	WeatherMatch  bool    `json:"weather_match"`  // 是否適合當前天氣
}

// WeatherInfo 天氣資訊
type WeatherInfo struct {
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
	Humidity    int     `json:"humidity"`
	Description string  `json:"description"`
}
