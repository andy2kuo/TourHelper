package frontend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andy2kuo/TourHelper/internal/bot/line"
	"github.com/andy2kuo/TourHelper/internal/bot/telegram"
	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/andy2kuo/TourHelper/internal/server"
	"github.com/gin-gonic/gin"
)

// HTTPServer HTTP 伺服器實作
type HTTPServer struct {
	router     *gin.Engine
	opt        *server.Options
	httpServer *http.Server
}

// Init 建立新的 HTTP 伺服器
func (s *HTTPServer) Init(opts *server.Options) error {
	// 設定 Gin 模式
	if opts.ServiceEnv == "release" {
		gin.SetMode(gin.ReleaseMode)
		logger.Info("Gin 設定為 Release 模式")
	} else {
		logger.Info("Gin 設定為 Debug 模式")
	}

	// 建立 Gin router
	r := gin.Default()

	s.router = r
	s.opt = opts

	// 註冊路由
	s.setupRoutes()

	return nil
}

// setupRoutes 設定所有路由
func (s *HTTPServer) setupRoutes() {
	// 設定靜態檔案服務（用於 Vue.js 前端）
	s.router.Static("/assets", "./web/dist/assets")
	s.router.StaticFile("/", "./web/dist/index.html")

	// Line Bot webhook
	if s.opt.Config.Line.Enabled {
		lineBot := line.NewBot(
			s.opt.Config.Line.ChannelSecret,
			s.opt.Config.Line.ChannelAccessToken,
		)
		s.router.POST("/webhook/line", lineBot.HandleWebhook)
		logger.Info("Line Bot 已啟用")
	}

	// Telegram Bot webhook
	if s.opt.Config.Telegram.Enabled {
		telegramBot := telegram.NewBot(s.opt.Config.Telegram.Token)
		s.router.POST("/webhook/telegram", telegramBot.HandleWebhook)
		logger.Info("Telegram Bot 已啟用")
	}

	// TODO: 在此處添加更多路由
}

// Start 啟動 HTTP 伺服器
func (s *HTTPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.opt.Config.Server.Host, s.opt.Config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Infof("HTTP 伺服器啟動於 %s", addr)

	// 檢查是否有配置 SSL 憑證
	if s.opt.Config.Server.CertFile != "" {
		// 使用 HTTPS
		logger.Infof("伺服器以 HTTPS 模式啟動於 %s", addr)
		logger.Infof("使用憑證檔案: %s", s.opt.Config.Server.CertFile)

		// 決定使用的私鑰檔案
		keyFile := s.opt.Config.Server.KeyFile
		if keyFile == "" {
			// 如果沒有指定 KeyFile，使用 CertFile（假設憑證和私鑰在同一個 PEM 檔案中）
			keyFile = s.opt.Config.Server.CertFile
			logger.Info("使用合併的 PEM 檔案（憑證和私鑰在同一檔案）")
		} else {
			logger.Infof("使用私鑰檔案: %s", keyFile)
		}

		// 啟動 HTTPS 伺服器
		if err := s.httpServer.ListenAndServeTLS(s.opt.Config.Server.CertFile, keyFile); err != nil && err != http.ErrServerClosed {
			return err
		}
	} else {
		// 使用 HTTP
		logger.Infof("伺服器以 HTTP 模式啟動於 %s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}
	}

	return nil
}

// Stop 停止 HTTP 伺服器
func (s *HTTPServer) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Name 返回伺服器名稱
func (s *HTTPServer) Name() string {
	return "HTTP Server"
}
