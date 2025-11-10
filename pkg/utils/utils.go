package utils

import (
	"math"
)

// CalculateDistance 計算兩個座標之間的距離（公里）
// 使用 Haversine 公式
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // 地球半徑（公里）

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

// EstimateTravelTime 估算旅行時間（分鐘）
// averageSpeed: 平均速度（km/h）
func EstimateTravelTime(distance, averageSpeed float64) int {
	if averageSpeed <= 0 {
		averageSpeed = 50 // 預設速度 50 km/h
	}
	return int(distance / averageSpeed * 60)
}

// Contains 檢查字串陣列是否包含特定字串
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
