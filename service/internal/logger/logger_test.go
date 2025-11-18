package logger

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/sirupsen/logrus"
)

func TestLoggerInit(t *testing.T) {
	// 測試用的設定
	logConfig := config.LogConfig{
		Level:      "debug",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	// 初始化 logger
	logger := Init("test_service", "dev", logConfig)

	if logger == nil {
		t.Fatal("Logger 初始化失敗")
	}

	// 檢查 logger 等級
	if logger.GetLevel() != logrus.DebugLevel {
		t.Errorf("期望 log level 為 DebugLevel，實際為 %v", logger.GetLevel())
	}

	// 測試寫入日誌
	logger.Info("測試日誌訊息")
	logger.Debug("除錯訊息")
	logger.WithField("key", "value").Info("帶欄位的日誌")

	// 清理
	defer func() {
		logger.Close()
		os.RemoveAll("./log")
	}()
}

func TestRotateWriter(t *testing.T) {
	// 建立臨時目錄
	tempDir := filepath.Join(os.TempDir(), "test_logs")
	defer os.RemoveAll(tempDir)

	// 建立輪換寫入器
	rw := NewRotateWriter(tempDir, "test", 1, 1, false)
	defer rw.Close()

	// 寫入測試資料
	testData := []byte("測試日誌資料\n")
	n, err := rw.Write(testData)
	if err != nil {
		t.Fatalf("寫入失敗: %v", err)
	}

	if n != len(testData) {
		t.Errorf("期望寫入 %d bytes，實際寫入 %d bytes", len(testData), n)
	}

	// 檢查檔案是否存在
	currentDate := time.Now().Format("2006-01-02")
	expectedFile := filepath.Join(tempDir, "test_"+currentDate+".log")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("期望的日誌檔案不存在: %s", expectedFile)
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// 初始化 logger
	logConfig := config.LogConfig{
		Level:      "info",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   false,
	}

	Init("test_convenience", "dev", logConfig)
	defer func() {
		GetLogger().Close()
		os.RemoveAll("./log")
	}()

	// 測試便利函數
	Info("測試 Info")
	Infof("測試 Infof: %s", "formatted")
	Debug("測試 Debug")
	Debugf("測試 Debugf: %d", 123)
	Warn("測試 Warn")
	Warnf("測試 Warnf: %v", true)
	Error("測試 Error")
	Errorf("測試 Errorf: %f", 3.14)

	// 測試 WithField 和 WithFields
	WithField("key", "value").Info("測試 WithField")
	WithFields(logrus.Fields{
		"key1": "value1",
		"key2": 123,
	}).Info("測試 WithFields")
}

func TestGetLoggerBeforeInit(t *testing.T) {
	// 儲存當前的 instance
	oldInstance := instance

	// 重置 instance 以測試未初始化的情況
	instance = nil

	// 恢復 instance，確保不影響其他測試
	defer func() {
		instance = oldInstance
		if r := recover(); r == nil {
			t.Error("期望 panic，但沒有發生")
		}
	}()

	// 呼叫 GetLogger 應該會 panic
	GetLogger()
}
