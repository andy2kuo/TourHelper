package lobby

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/andy2kuo/TourHelper/internal/server"
	"github.com/gin-gonic/gin"
)

// LobbyServer 大廳伺服器,使用 Gin 框架
// 主要功能:處理會員登入驗證,定時從 Redis 取得 Tour Server 狀態
type LobbyServer struct {
	router     *gin.Engine
	opt        *server.Options
	httpServer *http.Server
	// TODO: 新增 Redis 客戶端
	// redisClient *redis.Client
}

// Init 初始化伺服器
func (s *LobbyServer) Init(opts *server.Options) error {
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
func (s *LobbyServer) setupRoutes() {
	// 健康檢查路由
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": s.opt.ServiceName,
			"env":     s.opt.ServiceEnv,
			"version": s.opt.Version,
		})
	})

	// 登入相關路由
	auth := s.router.Group("/auth")
	{
		// 會員登入驗證
		auth.POST("/login", s.handleLogin)

		// TODO: LINE 第三方登入驗證功能待補
		// auth.POST("/line/callback", s.handleLineLogin)

		// 登出
		auth.POST("/logout", s.handleLogout)

		// 驗證 Token
		auth.POST("/verify", s.handleVerifyToken)
	}

	// 會員資訊路由
	member := s.router.Group("/member")
	{
		// 取得會員資訊
		member.GET("/:id", s.handleGetMemberInfo)

		// 更新會員資訊
		member.PUT("/:id", s.handleUpdateMemberInfo)
	}

	logger.Info("Lobby 路由已設定完成")
}

// Start 啟動 HTTP 伺服器
func (s *LobbyServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.opt.Config.Server.Host, s.opt.Config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Infof("Lobby 伺服器啟動於 %s", addr)

	// 檢查是否有配置 SSL 憑證
	if s.opt.Config.Server.CertFile != "" {
		// 使用 HTTPS
		logger.Infof("伺服器以 HTTPS 模式啟動於 %s", addr)
		if err := s.httpServer.ListenAndServeTLS(s.opt.Config.Server.CertFile, s.opt.Config.Server.KeyFile); err != nil && err != http.ErrServerClosed {
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
func (s *LobbyServer) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	logger.Info("正在關閉 Lobby 伺服器...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Name 返回伺服器名稱
func (s *LobbyServer) Name() string {
	return "Lobby Server"
}

// UpdateMemberStatusToRedis 更新會員狀態到 Redis
func (s *LobbyServer) UpdateMemberStatusToRedis(memberID string, status map[string]interface{}) error {
	// TODO: 實作 Redis 狀態更新
	// 1. 序列化狀態資料為 JSON
	// 2. 儲存到 Redis (key: member:{memberID}:status)
	// 3. 設定過期時間
	logger.WithFields(map[string]interface{}{
		"member_id": memberID,
		"status":    status,
	}).Info("TODO: 更新會員狀態到 Redis")
	return nil
}

// GetTourStatusFromRedis 從 Redis 取得 Tour Server 狀態
func (s *LobbyServer) GetTourStatusFromRedis() (map[string]interface{}, error) {
	// TODO: 實作從 Redis 讀取 Tour Server 狀態
	// 1. 從 Redis 讀取 key: tour:server:status
	// 2. 反序列化 JSON 為 map
	// 3. 回傳狀態資訊
	logger.Info("TODO: 從 Redis 取得 Tour Server 狀態")
	return map[string]interface{}{
		"status":  "unknown",
		"message": "Redis integration not implemented",
	}, nil
}

// StartTourStatusMonitor 啟動定時監控 Tour Server 狀態的協程
func (s *LobbyServer) StartTourStatusMonitor() {
	// TODO: 實作定時從 Redis 取得 Tour Server 狀態
	// 1. 建立 ticker,每 N 秒執行一次
	// 2. 呼叫 GetTourStatusFromRedis 取得狀態
	// 3. 更新本地快取或觸發相應處理
	logger.Info("TODO: 啟動 Tour Server 狀態監控協程")
}
