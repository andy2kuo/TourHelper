package services

import (
	"sync"

	"github.com/andy2kuo/TourHelper/internal/dao"
)

// Services 集中管理所有 service 實例
type Services struct {
	Recommendation RecommendationService
	Weather        WeatherService
	// 未來可以新增其他 service，例如：
	// User           UserService
	// Destination    DestinationService
}

var (
	instance *Services
	once     sync.Once
)

// Get 取得 services 實例（單例模式）
func Get() *Services {
	once.Do(func() {
		daos := dao.Get()
		instance = &Services{
			Recommendation: NewRecommendationService(daos),
			Weather:        NewWeatherService(),
			// 初始化其他 service
		}
	})
	return instance
}
