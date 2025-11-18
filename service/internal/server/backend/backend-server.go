package backend

import (
	"context"
	"fmt"
	"net/http"

	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/andy2kuo/TourHelper/internal/server"
	"github.com/gin-gonic/gin"
)

// BackendServer 後台管理伺服器,使用 Gin 框架
// 主要功能:後台管理功能,包含使用者登入驗證
type BackendServer struct {
	router     *gin.Engine
	opt        *server.Options
	httpServer *http.Server
	// TODO: 新增 Redis 客戶端
	// redisClient *redis.Client
}

// Init 初始化伺服器
func (s *BackendServer) Init(opts *server.Options) error {
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
func (s *BackendServer) setupRoutes() {
	// 健康檢查路由
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": s.opt.ServiceName,
			"env":     s.opt.ServiceEnv,
			"version": s.opt.Version,
		})
	})

	// 後台管理員登入驗證路由
	auth := s.router.Group("/admin/auth")
	{
		// TODO: 實作後台管理員登入驗證
		auth.POST("/login", s.handleAdminLogin)

		// TODO: 實作後台管理員登出
		auth.POST("/logout", s.handleAdminLogout)

		// TODO: 實作 Token 驗證
		auth.POST("/verify", s.handleVerifyToken)
	}

	// 會員管理路由 (需要驗證)
	// TODO: 實作中介層驗證 Token
	member := s.router.Group("/admin/member")
	// member.Use(s.authMiddleware()) // TODO: 實作驗證中介層
	{
		// TODO: 實作會員列表查詢
		member.GET("/list", s.handleGetMemberList)

		// TODO: 實作會員詳細資訊
		member.GET("/:id", s.handleGetMemberDetail)

		// TODO: 實作會員狀態管理(啟用/停用)
		member.PUT("/:id/status", s.handleUpdateMemberStatus)

		// TODO: 實作會員刪除
		member.DELETE("/:id", s.handleDeleteMember)
	}

	// Tour Server 管理路由 (需要驗證)
	tour := s.router.Group("/admin/tour")
	// tour.Use(s.authMiddleware()) // TODO: 實作驗證中介層
	{
		// TODO: 實作 Tour Server 狀態查詢
		tour.GET("/status", s.handleGetTourStatus)

		// TODO: 實作景點管理
		tour.GET("/destinations", s.handleGetDestinations)
		tour.POST("/destinations", s.handleCreateDestination)
		tour.PUT("/destinations/:id", s.handleUpdateDestination)
		tour.DELETE("/destinations/:id", s.handleDeleteDestination)
	}

	// 系統設定路由 (需要驗證)
	system := s.router.Group("/admin/system")
	// system.Use(s.authMiddleware()) // TODO: 實作驗證中介層
	{
		// TODO: 實作系統設定查詢
		system.GET("/config", s.handleGetSystemConfig)

		// TODO: 實作系統設定更新
		system.PUT("/config", s.handleUpdateSystemConfig)

		// TODO: 實作系統日誌查詢
		system.GET("/logs", s.handleGetSystemLogs)
	}

	logger.Info("Backend 路由已設定完成")
}

// Start 啟動 HTTP 伺服器
func (s *BackendServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.opt.Config.Server.Host, s.opt.Config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	logger.Infof("Backend 伺服器啟動於 %s", addr)

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
func (s *BackendServer) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	logger.Info("正在關閉 Backend 伺服器...")

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

// Name 返回伺服器名稱
func (s *BackendServer) Name() string {
	return "Backend Server"
}
