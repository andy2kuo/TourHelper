package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	instance *Logger
	once     sync.Once
)

// Logger 封裝 logrus.Logger 並提供額外功能
type Logger struct {
	*logrus.Logger
	rotateWriter *RotateWriter
	mu           sync.RWMutex
}

// RotateWriter 實作按日期和大小自動輪換的日誌寫入器
type RotateWriter struct {
	baseDir     string
	baseName    string
	currentDate string
	currentFile *lumberjack.Logger
	maxSize     int
	maxAge      int
	compress    bool
	mu          sync.Mutex
}

// getLogOutputPath 根據環境和作業系統返回日誌檔案的基礎路徑
func getLogOutputPath(serviceName, env string) string {
	if env == "dev" {
		// 開發環境：專案目錄下的 log 資料夾
		return "./log"
	}

	// 生產環境：根據作業系統選擇路徑
	if runtime.GOOS == "windows" {
		// Windows: C:/ProgramData/{SERVICE_NAME}/log
		return filepath.Join("C:", "ProgramData", serviceName, "log")
	}

	// Linux/Unix: /var/log/{SERVICE_NAME}
	return filepath.Join("/var/log", serviceName)
}

// NewRotateWriter 建立新的輪換寫入器
func NewRotateWriter(baseDir, baseName string, maxSize, maxAge int, compress bool) *RotateWriter {
	// 建立日誌目錄，權限設為 0755 (rwxr-xr-x)，所有用戶可讀可執行
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		logrus.WithError(err).Error("無法建立日誌目錄")
	}

	rw := &RotateWriter{
		baseDir:  baseDir,
		baseName: baseName,
		maxSize:  maxSize,
		maxAge:   maxAge,
		compress: compress,
	}

	rw.rotate()

	// 啟動定時清理任務
	go rw.cleanOldLogs()

	// 啟動定時檢查日期變更
	go rw.checkDateRotation()

	return rw
}

// Write 實作 io.Writer 介面
func (rw *RotateWriter) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	// 檢查是否需要按日期輪換
	currentDate := time.Now().Format("2006-01-02")
	if currentDate != rw.currentDate {
		rw.rotate()
	}

	return rw.currentFile.Write(p)
}

// rotate 執行日誌輪換
func (rw *RotateWriter) rotate() {
	currentDate := time.Now().Format("2006-01-02")
	filename := filepath.Join(rw.baseDir, fmt.Sprintf("%s_%s.log", rw.baseName, currentDate))

	if rw.currentFile != nil {
		rw.currentFile.Close()
	}

	// 在 Linux 系統上，臨時設定 umask 為 0022，確保檔案以 644 權限建立
	var oldMask int
	if runtime.GOOS != "windows" {
		oldMask = setUmask(0022)
		defer setUmask(oldMask)
	}

	rw.currentFile = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    rw.maxSize, // MB
		MaxBackups: 0,          // 不限制備份數量，由 cleanOldLogs 處理
		MaxAge:     0,          // 不使用內建的清理，由 cleanOldLogs 處理
		LocalTime:  true,
		Compress:   rw.compress,
	}

	// 確保日誌文件被建立後設定權限為 644 (rw-r--r--)，所有用戶可讀，擁有者可寫
	// 先嘗試寫入一個空字節以確保文件被建立
	rw.currentFile.Write([]byte(""))
	if runtime.GOOS != "windows" {
		if err := os.Chmod(filename, 0644); err != nil {
			logrus.WithError(err).WithField("file", filename).Warn("無法設定日誌文件權限")
		}
	}

	rw.currentDate = currentDate
}

// checkDateRotation 定時檢查日期是否變更，需要輪換日誌
func (rw *RotateWriter) checkDateRotation() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		currentDate := time.Now().Format("2006-01-02")
		rw.mu.Lock()
		if currentDate != rw.currentDate {
			rw.rotate()
		}
		rw.mu.Unlock()
	}
}

// cleanOldLogs 定時清理過期的日誌檔案
func (rw *RotateWriter) cleanOldLogs() {
	// 每小時檢查一次
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rw.doClean()
	}
}

// doClean 執行清理過期日誌
func (rw *RotateWriter) doClean() {
	cutoff := time.Now().AddDate(0, 0, -rw.maxAge)

	filepath.Walk(rw.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// 檢查是否為日誌檔案
		if filepath.Ext(path) == ".log" && info.ModTime().Before(cutoff) {
			os.Remove(path)
			logrus.WithField("file", path).Info("刪除過期日誌檔案")
		}

		return nil
	})
}

// Close 關閉日誌寫入器
func (rw *RotateWriter) Close() error {
	rw.mu.Lock()
	defer rw.mu.Unlock()

	if rw.currentFile != nil {
		return rw.currentFile.Close()
	}
	return nil
}

// Init 初始化 logger（單例模式）
func Init(serviceName, env string, logConfig config.LogConfig) *Logger {
	once.Do(func() {
		logDir := getLogOutputPath(serviceName, env)

		// 建立日誌目錄，權限設為 0755 (rwxr-xr-x)，所有用戶可讀可執行
		if err := os.MkdirAll(logDir, 0755); err != nil {
			fmt.Printf("無法建立日誌目錄 %s: %v\n", logDir, err)
			return
		}

		instance = &Logger{
			Logger: logrus.New(),
		}

		// 設定日誌等級
		logLevel := logrus.InfoLevel
		if env == "dev" {
			logLevel = logrus.DebugLevel
		}
		if logConfig.Level != "" {
			if level, err := logrus.ParseLevel(logConfig.Level); err == nil {
				logLevel = level
			}
		}
		instance.SetLevel(logLevel)

		// 設定格式（使用 JSON 格式，更適合日誌分析）
		instance.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})

		// 建立輪換寫入器
		instance.rotateWriter = NewRotateWriter(
			logDir,
			serviceName,
			logConfig.MaxSize,
			logConfig.MaxAge,
			logConfig.Compress,
		)

		// 設定輸出到檔案和控制台
		multiWriter := io.MultiWriter(os.Stdout, instance.rotateWriter)
		instance.SetOutput(multiWriter)

		instance.WithFields(logrus.Fields{
			"service":  serviceName,
			"env":      env,
			"logDir":   logDir,
			"maxSize":  logConfig.MaxSize,
			"maxAge":   logConfig.MaxAge,
			"compress": logConfig.Compress,
		}).Info("Logger 初始化完成")
	})

	return instance
}

// GetLogger 取得 logger 實例
func GetLogger() *Logger {
	if instance == nil {
		panic("Logger 尚未初始化，請先呼叫 Init()")
	}
	return instance
}

// Close 關閉 logger
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.rotateWriter != nil {
		return l.rotateWriter.Close()
	}
	return nil
}

// 便利函數：提供全域函數方便使用

// Debug 記錄 debug 級別日誌
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

// Debugf 記錄格式化的 debug 級別日誌
func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format, args...)
}

// Info 記錄 info 級別日誌
func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

// Infof 記錄格式化的 info 級別日誌
func Infof(format string, args ...interface{}) {
	GetLogger().Infof(format, args...)
}

// Warn 記錄 warn 級別日誌
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

// Warnf 記錄格式化的 warn 級別日誌
func Warnf(format string, args ...interface{}) {
	GetLogger().Warnf(format, args...)
}

// Error 記錄 error 級別日誌
func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

// Errorf 記錄格式化的 error 級別日誌
func Errorf(format string, args ...interface{}) {
	GetLogger().Errorf(format, args...)
}

// Fatal 記錄 fatal 級別日誌並結束程式
func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

// Fatalf 記錄格式化的 fatal 級別日誌並結束程式
func Fatalf(format string, args ...interface{}) {
	GetLogger().Fatalf(format, args...)
}

// WithField 建立帶有單一欄位的 logger entry
func WithField(key string, value interface{}) *logrus.Entry {
	return GetLogger().WithField(key, value)
}

// WithFields 建立帶有多個欄位的 logger entry
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetLogger().WithFields(fields)
}
