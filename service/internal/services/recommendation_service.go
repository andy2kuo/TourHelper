package services

import (
	"github.com/andy2kuo/TourHelper/internal/dao"
)

// RecommendationService 推薦服務介面
type RecommendationService interface {
	// TODO: 實作推薦業務邏輯
}

// recommendationService 推薦服務實作
type recommendationService struct {
	dao *dao.DAO
}

// NewRecommendationService 建立推薦服務
func NewRecommendationService(d *dao.DAO) RecommendationService {
	return &recommendationService{dao: d}
}
