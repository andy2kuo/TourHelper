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
	Log      LogConfig
}

// ServerConfig HTTP 伺服器設定
type ServerConfig struct {
	Host    string
	Port    string
	Env     string
	Version string
}

// DatabaseConfig 資料庫設定（MySQL Master-Slave 架構）
type DatabaseConfig struct {
	Masters []MasterDBConfig // Master 資料庫列表（可根據 Schema 區分）
	Slaves  []SlaveDBConfig  // Slave 資料庫列表（數量不定，可為空）

	// 全域連線池設定（所有 Master/Slave 共用）
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // 秒
}

// MasterDBConfig Master 資料庫設定
type MasterDBConfig struct {
	Name     string // Master 識別名稱（例如：main, analytics）
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
	ParseTime bool
	Loc      string

	// 可選：針對特定 Schema 的設定
	Schema string // 如果需要根據 Schema 區分 Master

	// 可選：此 Master 專用的連線池設定（覆蓋全域設定）
	MaxIdleConns    *int // nil 表示使用全域設定
	MaxOpenConns    *int
	ConnMaxLifetime *int
}

// SlaveDBConfig Slave 資料庫設定
type SlaveDBConfig struct {
	Name     string // Slave 識別名稱（例如：slave1, slave2）
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
	ParseTime bool
	Loc      string

	// 負載平衡權重（數值越大，被選中機率越高）
	Weight int

	// 可選：此 Slave 專用的連線池設定
	MaxIdleConns    *int
	MaxOpenConns    *int
	ConnMaxLifetime *int

	// 可選：指定此 Slave 對應的 Master
	MasterName string
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

// LogConfig 日誌設定
type LogConfig struct {
	Level      string // 日誌等級: debug, info, warn, error, fatal
	MaxSize    int    // MB
	MaxBackups int    // 保留的舊日誌檔案數量
	MaxAge     int    // 保留的天數
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

// DefaultConfig 返回預設設定
func DefaultConfig() *Config {
	cfg, err := Load("TourHelper", "dev", "0.0.1")
	if err != nil {
		panic(fmt.Sprintf("無法載入預設設定: %v", err))
	}
	return cfg
}

// setDefaults 設定預設值
func setDefaults() {
	// Server 預設值
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")

	// Database 預設值（MySQL Master-Slave）
	// 全域連線池設定
	viper.SetDefault("database.maxidleconns", 10)
	viper.SetDefault("database.maxopenconns", 100)
	viper.SetDefault("database.connmaxlifetime", 3600) // 1 小時

	// 預設 Master 設定（向後相容）
	viper.SetDefault("database.masters", []map[string]any{
		{
			"name":      "main",
			"host":      "localhost",
			"port":      "3306",
			"user":      "root",
			"password":  "",
			"dbname":    "tourhelper",
			"charset":   "utf8mb4",
			"parsetime": true,
			"loc":       "Local",
		},
	})

	// 預設 Slave 設定（空陣列，表示沒有 Slave）
	viper.SetDefault("database.slaves", []map[string]any{})

	// Line Bot 預設值
	viper.SetDefault("line.enabled", false)

	// Telegram Bot 預設值
	viper.SetDefault("telegram.enabled", false)

	// Weather 預設值
	viper.SetDefault("weather.provider", "openweathermap")

	// Maps 預設值
	viper.SetDefault("maps.provider", "google")

	// Log 預設值
	viper.SetDefault("log.maxsize", 100)   // 100 MB
	viper.SetDefault("log.maxbackups", 3)  // 保留 3 個備份
	viper.SetDefault("log.maxage", 28)     // 保留 28 天
	viper.SetDefault("log.compress", true) // 壓縮舊檔案
}
