package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/andy2kuo/TourHelper/internal/server"
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
	
	fmt.Println("設定檔載入成功")

	// 初始化 Logger（使用設定檔中的 log 設定）
	logger.Init(SERVICE_NAME, SERVICE_ENV, cfg.Log)

	
}

func main() {
	// 建立伺服器選項
	opts := &server.Options{
		Config:      cfg,
		ServiceName: SERVICE_NAME,
		ServiceEnv:  SERVICE_ENV,
		Version:     SERVICE_VERSION,
		Blocking:    false, // 非阻塞模式
	}

	// 建立 HTTP 伺服器
	httpServer := server.NewHTTPServer(opts)

	// 啟動伺服器
	if err := httpServer.Start(); err != nil {
		logger.Fatalf("伺服器啟動失敗: %v", err)
	}

	// 等待中斷信號以優雅地關閉伺服器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("收到關閉信號，正在關閉伺服器...")

	// 優雅關閉伺服器
	if err := httpServer.Stop(); err != nil {
		logger.Errorf("伺服器關閉錯誤: %v", err)
	}

	logger.Info("應用程式已關閉")
}
