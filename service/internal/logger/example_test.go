package logger_test

import (
	"fmt"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/sirupsen/logrus"
)

// ExampleInit 展示如何初始化和使用 logger
func ExampleInit() {
	// 設定 logger 配置
	logConfig := config.LogConfig{
		Level:      "info",
		MaxSize:    100, // 100 MB
		MaxBackups: 3,   // 保留 3 個備份
		MaxAge:     28,  // 保留 28 天
		Compress:   true,
	}

	// 初始化 logger
	// 注意：logger 使用單例模式，在實際應用中只需初始化一次
	log := logger.Init("example_service", "dev", logConfig)

	// 使用 logger 記錄日誌
	log.Info("應用程式啟動")

	// 使用便利函數
	logger.Info("這是一般訊息")
	logger.Warn("這是警告訊息")

	// 使用格式化字串
	logger.Infof("使用者 %s 已登入", "Alice")

	// 使用結構化欄位
	logger.WithField("user_id", 123).Info("使用者操作")
	logger.WithFields(logrus.Fields{
		"user_id": 123,
		"action":  "login",
	}).Info("使用者登入事件")

	fmt.Println("Logger 初始化和使用完成")

	// Output:
	// Logger 初始化和使用完成
}
