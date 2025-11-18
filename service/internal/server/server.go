package server

import (
	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/logger"
)

// Server 介面定義了所有伺服器類型需要實作的方法
type Server interface {
	// Init 初始化伺服器
	Init(opts Options) error

	// Start 啟動伺服器
	Start() error

	// Stop 停止伺服器
	Stop() error

	// Name 返回伺服器名稱
	Name() string
}

// Options 伺服器建立選項
type Options struct {
	Config      *config.Config // 應用程式設定
	ServiceName string         // 服務名稱
	ServiceEnv  string         // 服務環境
	Version     string         // 服務版本
	Blocking    bool           // 是否阻塞運行
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
		Blocking:    true,
	}
}

func StartServer(srv Server, opts *Options) error {
	if opts == nil {
		opts = defaultOpetion()

		logger.Warnf("%v 使用預設伺服器選項", srv.Name())
	}

	logger.WithFields(map[string]interface{}{
		"service": opts.ServiceName,
		"env":     opts.ServiceEnv,
		"version": opts.Version,
	}).Infof("%v 以 %v 模式啟動，版本 %v", opts.ServiceName, opts.ServiceEnv, opts.Version)

	if err := srv.Init(*opts); err != nil {
		return err
	}

	if opts.Blocking {
		logger.Infof("%v 以阻塞模式啟動", srv.Name())
		return srv.Start()
	} else {
		logger.Infof("%v 以非阻塞模式啟動", srv.Name())
		go srv.Start()
		return nil
	}
}
