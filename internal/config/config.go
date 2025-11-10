package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 包含所有應用程式設定
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Line     LineBotConfig
	Telegram TelegramBotConfig
	Weather  WeatherConfig
	Maps     MapsConfig
}

// ServerConfig HTTP 伺服器設定
type ServerConfig struct {
	Host string
	Port string
	Mode string // debug, release
}

// DatabaseConfig 資料庫設定
type DatabaseConfig struct {
	Type     string // mysql, postgres, sqlite
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LineBotConfig Line Bot 設定
type LineBotConfig struct {
	Enabled            bool
	ChannelSecret      string
	ChannelAccessToken string
}

// TelegramBotConfig Telegram Bot 設定
type TelegramBotConfig struct {
	Enabled bool
	Token   string
}

// WeatherConfig 天氣 API 設定
type WeatherConfig struct {
	APIKey   string
	Provider string // openweathermap, weatherapi, etc.
}

// MapsConfig 地圖 API 設定
type MapsConfig struct {
	APIKey   string
	Provider string // google, here, mapbox, etc.
}

// Load 載入設定檔
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// 設定環境變數前綴
	viper.SetEnvPrefix("TOURHELPER")
	viper.AutomaticEnv()

	// 設定預設值
	setDefaults()

	// 讀取設定檔
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 設定檔不存在，使用預設值和環境變數
			fmt.Println("警告: 找不到設定檔，使用預設值")
		} else {
			return nil, fmt.Errorf("讀取設定檔錯誤: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析設定檔錯誤: %w", err)
	}

	return &config, nil
}

// setDefaults 設定預設值
func setDefaults() {
	// Server 預設值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// Database 預設值
	viper.SetDefault("database.type", "sqlite")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.dbname", "tourhelper.db")
	viper.SetDefault("database.sslmode", "disable")

	// Line Bot 預設值
	viper.SetDefault("line.enabled", false)

	// Telegram Bot 預設值
	viper.SetDefault("telegram.enabled", false)

	// Weather 預設值
	viper.SetDefault("weather.provider", "openweathermap")

	// Maps 預設值
	viper.SetDefault("maps.provider", "google")
}
