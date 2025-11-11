package utils

import (
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name     string
		lat1     float64
		lon1     float64
		lat2     float64
		lon2     float64
		expected float64
		delta    float64 // 允許的誤差範圍
	}{
		{
			name:     "台北101到陽明山",
			lat1:     25.0340,
			lon1:     121.5645,
			lat2:     25.1896,
			lon2:     121.5453,
			expected: 17.5, // 約17.5公里
			delta:    1.0,
		},
		{
			name:     "相同位置",
			lat1:     25.0340,
			lon1:     121.5645,
			lat2:     25.0340,
			lon2:     121.5645,
			expected: 0.0,
			delta:    0.1,
		},
		{
			name:     "台北到高雄",
			lat1:     25.0330,
			lon1:     121.5654,
			lat2:     22.6273,
			lon2:     120.3014,
			expected: 297.0, // 約297公里
			delta:    5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)
			diff := math.Abs(result - tt.expected)
			if diff > tt.delta {
				t.Errorf("CalculateDistance() = %v, 期望 %v (±%v)", result, tt.expected, tt.delta)
			}
		})
	}
}

func TestEstimateTravelTime(t *testing.T) {
	tests := []struct {
		name         string
		distance     float64
		averageSpeed float64
		expected     int
	}{
		{
			name:         "50公里，時速50",
			distance:     50.0,
			averageSpeed: 50.0,
			expected:     60, // 60分鐘
		},
		{
			name:         "100公里，時速100",
			distance:     100.0,
			averageSpeed: 100.0,
			expected:     60,
		},
		{
			name:         "使用預設速度",
			distance:     50.0,
			averageSpeed: 0, // 應使用預設50km/h
			expected:     60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EstimateTravelTime(tt.distance, tt.averageSpeed)
			if result != tt.expected {
				t.Errorf("EstimateTravelTime() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		item     string
		expected bool
	}{
		{
			name:     "找到項目",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "banana",
			expected: true,
		},
		{
			name:     "找不到項目",
			slice:    []string{"apple", "banana", "cherry"},
			item:     "orange",
			expected: false,
		},
		{
			name:     "空陣列",
			slice:    []string{},
			item:     "apple",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.slice, tt.item)
			if result != tt.expected {
				t.Errorf("Contains() = %v, 期望 %v", result, tt.expected)
			}
		})
	}
}
