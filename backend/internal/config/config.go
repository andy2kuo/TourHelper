package config

import (
	"fmt"
	"path/filepath"
	"runtime"

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
	Log      LogConfig
}

// ServerConfig HTTP 伺服器設定
type ServerConfig struct {
	Host    string
	Port    string
	Env     string
	Version string
}

// DatabaseConfig 資料庫設定（MySQL）
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // 秒
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

// LogConfig 日誌設定
type LogConfig struct {
	MaxSize    int // MB
	MaxBackups int // 保留的舊日誌檔案數量
	MaxAge     int // 保留的天數
	Compress   bool
}

// Load 載入設定檔
func Load(serviceName, env, version string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 根據環境和作業系統設定設定檔路徑
	if env == "dev" {
		// 開發環境：從專案目錄讀取
		viper.AddConfigPath("./configs")
		viper.AddConfigPath(".")
	} else {
		// 生產環境：根據作業系統選擇路徑
		var configPath string
		if runtime.GOOS == "windows" {
			// Windows: C:/ProgramData/{SERVICE_NAME}/configs
			configPath = filepath.Join("C:", "ProgramData", serviceName, "configs")
		} else {
			// Linux/Unix: /etc/{SERVICE_NAME}/configs
			configPath = filepath.Join("/etc", serviceName, "configs")
		}
		viper.AddConfigPath(configPath)
	}

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

	// 設定從啟動參數傳入的值
	config.Server.Env = env
	config.Server.Version = version

	return &config, nil
}

// setDefaults 設定預設值
func setDefaults() {
	// Server 預設值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// Database 預設值（MySQL）
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.dbname", "tourhelper")
	viper.SetDefault("database.charset", "utf8mb4")
	viper.SetDefault("database.parsetime", true)
	viper.SetDefault("database.loc", "Local")
	viper.SetDefault("database.maxidleconns", 10)
	viper.SetDefault("database.maxopenconns", 100)
	viper.SetDefault("database.connmaxlifetime", 3600) // 1 小時

	// Line Bot 預設值
	viper.SetDefault("line.enabled", false)

	// Telegram Bot 預設值
	viper.SetDefault("telegram.enabled", false)

	// Weather 預設值
	viper.SetDefault("weather.provider", "openweathermap")

	// Maps 預設值
	viper.SetDefault("maps.provider", "google")

	// Log 預設值
	viper.SetDefault("log.maxsize", 100)    // 100 MB
	viper.SetDefault("log.maxbackups", 3)   // 保留 3 個備份
	viper.SetDefault("log.maxage", 28)      // 保留 28 天
	viper.SetDefault("log.compress", true)  // 壓縮舊檔案
}
