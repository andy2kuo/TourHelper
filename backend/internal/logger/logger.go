package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// Init 初始化 logger
func Init(serviceName, env string, logConfig config.LogConfig) error {
	Log = logrus.New()

	// 設定 log 格式
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 設定 log level
	if env == "dev" {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}

	// 根據環境設定 log 路徑
	var logPath string
	if env == "dev" {
		// 開發環境：專案目錄下的 log 資料夾
		logPath = "./log"
	} else {
		// 生產環境：根據作業系統選擇路徑
		if runtime.GOOS == "windows" {
			// Windows: C:/ProgramData/{SERVICE_NAME}/log
			logPath = filepath.Join("C:", "ProgramData", serviceName, "log")
		} else {
			// Linux/Unix: /var/log/{SERVICE_NAME}
			logPath = filepath.Join("/var/log", serviceName)
		}
	}

	// 確保 log 目錄存在
	if err := os.MkdirAll(logPath, 0755); err != nil {
		return fmt.Errorf("無法建立 log 目錄: %w", err)
	}

	// 建立 log 檔案（使用服務名稱）
	logFile := filepath.Join(logPath, serviceName+".log")

	// 使用 lumberjack 實現 log rotation
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    logConfig.MaxSize,    // MB
		MaxBackups: logConfig.MaxBackups, // 保留的舊檔案數量
		MaxAge:     logConfig.MaxAge,     // 保留天數
		Compress:   logConfig.Compress,   // 是否壓縮
	}

	// 同時輸出到檔案和 console
	mw := io.MultiWriter(os.Stdout, fileWriter)
	Log.SetOutput(mw)

	Log.WithFields(logrus.Fields{
		"service":    serviceName,
		"env":        env,
		"logPath":    logFile,
		"maxSize":    logConfig.MaxSize,
		"maxBackups": logConfig.MaxBackups,
		"maxAge":     logConfig.MaxAge,
	}).Info("Logger 初始化完成")

	return nil
}

// Debug 記錄 debug 級別日誌
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf 記錄格式化的 debug 級別日誌
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Info 記錄 info 級別日誌
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof 記錄格式化的 info 級別日誌
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Warn 記錄 warn 級別日誌
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf 記錄格式化的 warn 級別日誌
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error 記錄 error 級別日誌
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf 記錄格式化的 error 級別日誌
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal 記錄 fatal 級別日誌並結束程式
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf 記錄格式化的 fatal 級別日誌並結束程式
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// WithFields 建立帶有欄位的 logger
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
