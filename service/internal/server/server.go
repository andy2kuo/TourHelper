package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/database"
	"github.com/andy2kuo/TourHelper/internal/logger"
)

// Server 介面定義了所有伺服器類型需要實作的方法
type Server interface {
	// Init 初始化伺服器
	Init(opts *Options) error

	// Start 啟動伺服器
	Start() error

	// Stop 停止伺服器
	Stop(context.Context) error

	// Name 返回伺服器名稱
	Name() string
}

// Options 伺服器建立選項
type Options struct {
	Config      *config.Config // 應用程式設定
	ServiceName string         // 服務名稱
	ServiceEnv  string         // 服務環境
	Version     string         // 服務版本
}

var envFormat = map[string]string{
	"dev":     "開發環境",
	"staging": "測試環境",
	"release": "正式環境",
}

var defaultOpetion = func() *Options {
	service_name := "default_service"
	service_env := "dev"
	service_version := "0.0.1-dev"

	return &Options{
		Config:      config.DefaultConfig(),
		ServiceName: service_name,
		ServiceEnv:  service_env,
		Version:     service_version,
	}
}

// GetEnvDescription 根據環境代碼返回描述字串
func GetEnvDescription(env string) (string, bool) {
	env = strings.ToLower(env)
	desc, ok := envFormat[env]

	return desc, ok
}

func StartServer(srv Server, opts *Options) error {
	if opts == nil {
		opts = defaultOpetion()

		logger.Warnf("%v 使用預設伺服器選項", srv.Name())
	}

	var envDesc string
	var envDescOk bool
	if envDesc, envDescOk = envFormat[opts.ServiceEnv]; !envDescOk {
		return fmt.Errorf("%v 不支援的服務環境: %v", srv.Name(), opts.ServiceEnv)
	}

	logger.WithFields(map[string]interface{}{
		"service": opts.ServiceName,
		"env":     opts.ServiceEnv,
		"version": opts.Version,
	}).Infof("%v 以 %v 啟動，版本 %v", opts.ServiceName, envDesc, opts.Version)

	// 初始化資料庫（MySQL 和 Redis）
	logger.Info("初始化資料庫連線...")
	if err := database.Init(opts.Config); err != nil {
		return fmt.Errorf("資料庫初始化失敗: %w", err)
	}

	if err := srv.Init(opts); err != nil {
		return err
	}

	go func() {
		err := srv.Start()
		if err != nil {
			logger.Errorf("%v 啟動失敗: %v", srv.Name(), err)
		}
	}()

	waitForShutdown(srv)

	return nil
}

func waitForShutdown(srv Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在關閉伺服器...")

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := srv.Stop(ctx); err != nil {
		logger.Errorf("伺服器關閉失敗: %v", err)
		return
	}

	// 關閉資料庫連線
	logger.Info("正在關閉資料庫連線...")
	if err := database.Close(); err != nil {
		logger.Errorf("資料庫關閉失敗: %v", err)
	}

	logger.Info("伺服器已關閉")
}
