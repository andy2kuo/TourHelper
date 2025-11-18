package main

import (
	"fmt"
	"strings"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/andy2kuo/TourHelper/internal/server"
	"github.com/andy2kuo/TourHelper/internal/server/frontend"
)

var SERVICE_NAME = "tour_helper"  // 預設值，會在編譯時透過 -ldflags 覆寫
var SERVICE_ENV = "dev"           // 預設值，會在編譯時透過 -ldflags 覆寫
var SERVICE_VERSION = "0.0.1-dev" // 預設值，會在編譯時透過 -ldflags 覆寫

var cfg *config.Config

func init() {
	SERVICE_ENV = strings.ToLower(SERVICE_ENV)

	var err error
	// 載入設定（需要先載入才能取得 log 設定）
	cfg, err = config.Load(SERVICE_NAME, SERVICE_ENV, SERVICE_VERSION)
	if err != nil {
		panic(fmt.Errorf("無法載入設定: %v", err))
	}

	// 初始化 Logger（使用設定檔中的 log 設定）
	err = logger.Init(SERVICE_NAME, SERVICE_ENV, cfg.Log)
	if err != nil {
		panic(fmt.Errorf("無法初始化 Logger: %v", err))
	}
}

func main() {
	defer logger.GetLogger().Close()

	// 建立伺服器選項
	opts := &server.Options{
		Config:      cfg,
		ServiceName: SERVICE_NAME,
		ServiceEnv:  SERVICE_ENV,
		Version:     SERVICE_VERSION,
	}

	logger.Infof("啟動應用程式 %s，版本 %s", SERVICE_NAME, SERVICE_VERSION)
	// 這裡會Block直到伺服器停止
	err := server.StartServer(&frontend.HTTPServer{}, opts)
	if err != nil {
		logger.Fatalf("伺服器啟動失敗: %v", err)
	}

	logger.Info("應用程式已關閉")
}
