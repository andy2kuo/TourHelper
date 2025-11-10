package models

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 資料庫設定
type DatabaseConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// InitDB 初始化資料庫連線
func InitDB(config DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	switch config.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.DBName), gormConfig)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.User, config.Password, config.Host, config.Port, config.DBName)
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		return nil, fmt.Errorf("不支援的資料庫類型: %s", config.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("連線資料庫失敗: %w", err)
	}

	// 自動遷移資料表結構
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("資料表遷移失敗: %w", err)
	}

	log.Println("資料庫連線成功")
	return db, nil
}

// autoMigrate 自動遷移所有資料表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&UserPreferences{},
		&Destination{},
		&Tag{},
		&SearchHistory{},
		&WeatherData{},
	)
}

// SeedSampleData 填充範例資料（開發用）
func SeedSampleData(db *gorm.DB) error {
	// 檢查是否已有資料
	var count int64
	db.Model(&Destination{}).Count(&count)
	if count > 0 {
		log.Println("資料庫已有資料，跳過範例資料填充")
		return nil
	}

	// 建立範例標籤
	tags := []Tag{
		{Name: "自然"},
		{Name: "文化"},
		{Name: "美食"},
		{Name: "購物"},
		{Name: "冒險"},
		{Name: "親子"},
		{Name: "網美"},
	}

	for i := range tags {
		if err := db.Create(&tags[i]).Error; err != nil {
			return err
		}
	}

	// 建立範例景點
	destinations := []Destination{
		{
			Name:        "陽明山國家公園",
			Description: "台北近郊的自然景觀勝地，擁有火山地形、溫泉、櫻花等豐富資源",
			Category:    "nature",
			Latitude:    25.1896,
			Longitude:   121.5453,
			Rating:      4.5,
			Address:     "台北市北投區竹子湖路",
			City:        "台北市",
			Region:      "北部",
			Country:     "Taiwan",
			Tags:        []Tag{tags[0], tags[5]}, // 自然、親子
		},
		{
			Name:        "淡水老街",
			Description: "歷史悠久的老街，有許多特色小吃和古蹟",
			Category:    "culture",
			Latitude:    25.1677,
			Longitude:   121.4406,
			Rating:      4.3,
			Address:     "新北市淡水區中正路",
			City:        "新北市",
			Region:      "北部",
			Country:     "Taiwan",
			Tags:        []Tag{tags[1], tags[2]}, // 文化、美食
		},
		{
			Name:        "九份老街",
			Description: "山城特色老街，以夜景和芋圓聞名",
			Category:    "culture",
			Latitude:    25.1095,
			Longitude:   121.8456,
			Rating:      4.4,
			Address:     "新北市瑞芳區基山街",
			City:        "新北市",
			Region:      "北部",
			Country:     "Taiwan",
			Tags:        []Tag{tags[1], tags[2], tags[6]}, // 文化、美食、網美
		},
		{
			Name:        "台北101",
			Description: "台北地標，擁有觀景台和購物中心",
			Category:    "shopping",
			Latitude:    25.0340,
			Longitude:   121.5645,
			Rating:      4.6,
			Address:     "台北市信義區信義路五段7號",
			City:        "台北市",
			Region:      "北部",
			Country:     "Taiwan",
			Tags:        []Tag{tags[3], tags[6]}, // 購物、網美
		},
		{
			Name:        "野柳地質公園",
			Description: "世界級地質景觀，以女王頭聞名",
			Category:    "nature",
			Latitude:    25.2055,
			Longitude:   121.6898,
			Rating:      4.5,
			Address:     "新北市萬里區野柳里港東路167-1號",
			City:        "新北市",
			Region:      "北部",
			Country:     "Taiwan",
			Tags:        []Tag{tags[0], tags[5], tags[6]}, // 自然、親子、網美
		},
	}

	for i := range destinations {
		if err := db.Create(&destinations[i]).Error; err != nil {
			return err
		}
	}

	log.Printf("已填充 %d 個範例景點", len(destinations))
	return nil
}
