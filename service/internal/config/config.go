package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 包含所有應用程式設定
type Config struct {
	Server   ServerConfig      `mapstructure:"server" json:"server" yaml:"server"`
	Database DatabaseConfig    `mapstructure:"database" json:"database" yaml:"database"`
	Line     LineBotConfig     `mapstructure:"line" json:"line" yaml:"line"`
	Telegram TelegramBotConfig `mapstructure:"telegram" json:"telegram" yaml:"telegram"`
	Log      LogConfig         `mapstructure:"log" json:"log" yaml:"log"`
}

// ServerConfig HTTP 伺服器設定
type ServerConfig struct {
	Host         string        `mapstructure:"host" json:"host" yaml:"host"`
	Port         int           `mapstructure:"port" json:"port" yaml:"port"`
	CertFile     string        `mapstructure:"certFile" json:"certFile" yaml:"certFile"`             // SSL 憑證檔案路徑 (.pem 或 .crt)
	KeyFile      string        `mapstructure:"keyFile" json:"keyFile" yaml:"keyFile"`                // SSL 私鑰檔案路徑 (.key)，如果憑證和私鑰在同一個 PEM 檔案中則不需要
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`    // HTTP 讀取超時時間
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"` // HTTP 寫入超時時間
	CORS         CORSConfig    `mapstructure:"cors" json:"cors" yaml:"cors"`
}

type CORSConfig struct {
	Enabled      bool     `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	AllowOrigins []string `mapstructure:"allowOrigins" json:"allowOrigins" yaml:"allowOrigins"`
	AllowMethods []string `mapstructure:"allowMethods" json:"allowMethods" yaml:"allowMethods"`
	AllowHeaders []string `mapstructure:"allowHeaders" json:"allowHeaders" yaml:"allowHeaders"`
}

// DatabaseConfig 資料庫設定（MySQL Master-Slave 架構）
type DatabaseConfig struct {
	Masters []MasterDBConfig `mapstructure:"masters" json:"masters" yaml:"masters"` // Master 資料庫列表（可根據 Database 區分）
	Slaves  []SlaveDBConfig  `mapstructure:"slaves" json:"slaves" yaml:"slaves"`    // Slave 資料庫列表（數量不定，可為空）

	// 全域連線池設定（所有 Master/Slave 共用）
	MaxIdleConns    int           `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime" yaml:"connMaxLifetime"` // 連線最大存活時間

	// SQL 日誌設定
	LogSlowQuery       bool          `mapstructure:"logSlowQuery" json:"logSlowQuery" yaml:"logSlowQuery"`                   // 是否記錄慢查詢
	SlowQueryThreshold time.Duration `mapstructure:"slowQueryThreshold" json:"slowQueryThreshold" yaml:"slowQueryThreshold"` // 慢查詢門檻
	LogAllQueries      bool          `mapstructure:"logAllQueries" json:"logAllQueries" yaml:"logAllQueries"`                // 是否記錄所有查詢（開發用）
}

// MasterDBConfig Master 資料庫設定
type MasterDBConfig struct {
	Name      string `mapstructure:"name" json:"name" yaml:"name"` // Master 識別名稱（例如：main, analytics）
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      string `mapstructure:"port" json:"port" yaml:"port"`
	User      string `mapstructure:"user" json:"user" yaml:"user"`
	Password  string `mapstructure:"password" json:"password" yaml:"password"`
	DBName    string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Charset   string `mapstructure:"charset" json:"charset" yaml:"charset"`
	ParseTime bool   `mapstructure:"parsetime" json:"parsetime" yaml:"parsetime"`
	Loc       string `mapstructure:"loc" json:"loc" yaml:"loc"`

	// 可選：針對特定 Database 的設定
	Database string `mapstructure:"database" json:"database" yaml:"database"` // 如果需要根據 Database 區分 Master

	// 可選：此 Master 專用的連線池設定（覆蓋全域設定）
	MaxIdleConns    *int           `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"` // nil 表示使用全域設定
	MaxOpenConns    *int           `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLifetime *time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime" yaml:"connMaxLifetime"`
}

// SlaveDBConfig Slave 資料庫設定
type SlaveDBConfig struct {
	Name      string `mapstructure:"name" json:"name" yaml:"name"` // Slave 識別名稱（例如：slave1, slave2）
	Host      string `mapstructure:"host" json:"host" yaml:"host"`
	Port      string `mapstructure:"port" json:"port" yaml:"port"`
	User      string `mapstructure:"user" json:"user" yaml:"user"`
	Password  string `mapstructure:"password" json:"password" yaml:"password"`
	DBName    string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Charset   string `mapstructure:"charset" json:"charset" yaml:"charset"`
	ParseTime bool   `mapstructure:"parsetime" json:"parsetime" yaml:"parsetime"`
	Loc       string `mapstructure:"loc" json:"loc" yaml:"loc"`

	// 負載平衡權重（數值越大，被選中機率越高）
	Weight int `mapstructure:"weight" json:"weight" yaml:"weight"`

	// 可選：此 Slave 專用的連線池設定
	MaxIdleConns    *int           `mapstructure:"maxIdleConns" json:"maxIdleConns" yaml:"maxIdleConns"`
	MaxOpenConns    *int           `mapstructure:"maxOpenConns" json:"maxOpenConns" yaml:"maxOpenConns"`
	ConnMaxLifetime *time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime" yaml:"connMaxLifetime"`

	// 可選：指定此 Slave 對應的 Master
	MasterName string `mapstructure:"masterName" json:"masterName" yaml:"masterName"`
}

// LineBotConfig Line Bot 設定
type LineBotConfig struct {
	Enabled            bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	ChannelSecret      string `mapstructure:"channelSecret" json:"channelSecret" yaml:"channelSecret"`
	ChannelAccessToken string `mapstructure:"channelAccessToken" json:"channelAccessToken" yaml:"channelAccessToken"`
}

// TelegramBotConfig Telegram Bot 設定
type TelegramBotConfig struct {
	Enabled bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	Token   string `mapstructure:"token" json:"token" yaml:"token"`
}

// LogConfig 日誌設定
type LogConfig struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`                // 日誌等級: debug, info, warn, error, fatal
	MaxSize    int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`          // MB
	MaxBackups int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"` // 保留的舊日誌檔案數量
	MaxAge     int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`             // 保留的天數
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
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
	viper.SetEnvPrefix(strings.ToUpper(serviceName))
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

	fmt.Println("設定檔載入成功")
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

	// Database 預設值（MySQL Master-Slave）
	// 全域連線池設定
	viper.SetDefault("database.maxidleconns", 10)
	viper.SetDefault("database.maxopenconns", 100)
	viper.SetDefault("database.connmaxlifetime", 3600) // 1 小時

	// SQL 日誌設定
	viper.SetDefault("database.logslowquery", true)                       // 預設啟用慢查詢記錄
	viper.SetDefault("database.slowquerythreshold", 200*time.Millisecond) // 預設 200ms 為慢查詢
	viper.SetDefault("database.logallqueries", false)                     // 預設不記錄所有查詢

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
	viper.SetDefault("log.maxSize", 100)   // 100 MB
	viper.SetDefault("log.maxBackups", 3)  // 保留 3 個備份
	viper.SetDefault("log.maxAge", 28)     // 保留 28 天
	viper.SetDefault("log.compress", true) // 壓縮舊檔案
}
