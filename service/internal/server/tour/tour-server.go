package tour

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

// TourServer 旅遊伺服器,使用 Gin 框架,提供 HTTP + WebSocket
// 主要功能:處理旅遊相關的業務邏輯
type TourServer struct {
	router     *gin.Engine
	opt        *server.Options
	httpServer *http.Server
	wsHub      *Hub // WebSocket Hub
}

// Init 初始化伺服器
func (s *TourServer) Init(opts *server.Options) error {
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

	// 建立並啟動 WebSocket Hub
	s.wsHub = NewHub()
	go s.wsHub.Run()
	logger.Info("WebSocket Hub 已啟動")

	// 註冊路由
	s.setupRoutes()

	return nil
}

// setupRoutes 設定所有路由
func (s *TourServer) setupRoutes() {
	// 設定靜態檔案服務（用於 Vue.js 前端）
	s.router.Static("/assets", "./web/dist/assets")
	s.router.StaticFile("/", "./web/dist/index.html")

	// WebSocket 路由
	wsHandler := NewWebSocketHandler(s.wsHub)
	s.router.GET("/ws", wsHandler.HandleWebSocket)
	s.router.GET("/ws/info", wsHandler.HandleWebSocketInfo)
	logger.Info("WebSocket 路由已設定: /ws")

	// 健康檢查路由
	healthCheck := NewHealthCheckHandler(s.opt.ServiceName, s.opt.ServiceEnv, s.opt.Version)
	s.router.GET("/health", healthCheck.Handle)

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

	// TODO: 在此處添加旅遊相關的 API 路由
	// 例如: 推薦景點、查詢景點資訊、使用者偏好設定等
}

// Start 啟動 HTTP 伺服器
func (s *TourServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.opt.Config.Server.Host, s.opt.Config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Infof("Tour 伺服器啟動於 %s", addr)

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
func (s *TourServer) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	logger.Info("正在關閉 Tour 伺服器...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Name 返回伺服器名稱
func (s *TourServer) Name() string {
	return "Tour Server"
}

// GetHub 取得 WebSocket Hub (供 Lobby Server 使用)
func (s *TourServer) GetHub() *Hub {
	return s.wsHub
}
