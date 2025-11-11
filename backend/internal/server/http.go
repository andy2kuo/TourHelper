package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andy2kuo/TourHelper/internal/bot/line"
	"github.com/andy2kuo/TourHelper/internal/bot/telegram"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/gin-gonic/gin"
)

// HTTPServer HTTP 伺服器實作
type HTTPServer struct {
	router      *gin.Engine
	server      *http.Server
	config      *Options
	httpServer  *http.Server
}

// NewHTTPServer 建立新的 HTTP 伺服器
func NewHTTPServer(opts *Options) *HTTPServer {
	// 設定 Gin 模式
	if opts.ServiceEnv == "release" || opts.ServiceEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Gin 設定為 Release 模式")
	} else {
		logger.Info("Gin 設定為 Debug 模式")
	}

	// 建立 Gin router
	r := gin.Default()

	s := &HTTPServer{
		router: r,
		config: opts,
	}

	// 註冊路由
	s.setupRoutes()

	return s
}

// setupRoutes 設定所有路由
func (s *HTTPServer) setupRoutes() {
	// 設定靜態檔案服務（用於 Vue.js 前端）
	s.router.Static("/assets", "./web/dist/assets")
	s.router.StaticFile("/", "./web/dist/index.html")

	// Line Bot webhook
	if s.config.Config.Line.Enabled {
		lineBot := line.NewBot(
			s.config.Config.Line.ChannelSecret,
			s.config.Config.Line.ChannelAccessToken,
		)
		s.router.POST("/webhook/line", lineBot.HandleWebhook)
		logger.Info("Line Bot 已啟用")
	}

	// Telegram Bot webhook
	if s.config.Config.Telegram.Enabled {
		telegramBot := telegram.NewBot(s.config.Config.Telegram.Token)
		s.router.POST("/webhook/telegram", telegramBot.HandleWebhook)
		logger.Info("Telegram Bot 已啟用")
	}

	// 健康檢查
	s.router.GET("/health", s.healthCheck)
}

// healthCheck 健康檢查端點
func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"service": s.config.ServiceName,
		"env":     s.config.Config.Server.Env,
		"version": s.config.Config.Server.Version,
	})
}

// Start 啟動 HTTP 伺服器
func (s *HTTPServer) Start() error {
	addr := s.config.Config.Server.Host + ":" + s.config.Config.Server.Port

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Infof("HTTP 伺服器啟動於 %s", addr)

	// 在 goroutine 中啟動，以便支援優雅關閉
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP 伺服器啟動失敗: %v", err)
		}
	}()

	return nil
}

// Stop 停止 HTTP 伺服器
func (s *HTTPServer) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	logger.Info("正在關閉 HTTP 伺服器...")

	// 設定 5 秒的超時時間來關閉伺服器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP 伺服器關閉失敗: %w", err)
	}

	logger.Info("HTTP 伺服器已關閉")
	return nil
}

// Name 返回伺服器名稱
func (s *HTTPServer) Name() string {
	return "HTTP Server"
}
