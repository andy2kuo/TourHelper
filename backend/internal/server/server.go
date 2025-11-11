package server

import (
	"github.com/andy2kuo/TourHelper/internal/config"
)

// Server 介面定義了所有伺服器類型需要實作的方法
type Server interface {
	// Start 啟動伺服器
	Start() error

	// Stop 停止伺服器
	Stop() error

	// Name 返回伺服器名稱
	Name() string
}

// Options 伺服器建立選項
type Options struct {
	Config      *config.Config
	ServiceName string
	ServiceEnv  string
	Version     string
}
