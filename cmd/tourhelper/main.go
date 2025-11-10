package main

import (
	"log"

	"github.com/andy2kuo/TourHelper/internal/config"
	"github.com/andy2kuo/TourHelper/internal/handlers"
	"github.com/andy2kuo/TourHelper/internal/bot/line"
	"github.com/andy2kuo/TourHelper/internal/bot/telegram"
	"github.com/gin-gonic/gin"
)

func main() {
	// 載入設定
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("無法載入設定: %v", err)
	}

	// 設定 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 建立 Gin router
	r := gin.Default()

	// 設定靜態檔案服務（用於 Vue.js 前端）
	r.Static("/assets", "./web/dist/assets")
	r.StaticFile("/", "./web/dist/index.html")

	// API 路由群組
	api := r.Group("/api/v1")
	{
		// 旅遊推薦相關路由
		api.GET("/recommendations", handlers.GetRecommendations)
		api.POST("/user/preferences", handlers.UpdatePreferences)
		api.GET("/user/preferences", handlers.GetPreferences)
	}

	// Line Bot webhook
	if cfg.Line.Enabled {
		lineBot := line.NewBot(cfg.Line.ChannelSecret, cfg.Line.ChannelAccessToken)
		r.POST("/webhook/line", lineBot.HandleWebhook)
		log.Println("Line Bot 已啟用")
	}

	// Telegram Bot webhook
	if cfg.Telegram.Enabled {
		telegramBot := telegram.NewBot(cfg.Telegram.Token)
		r.POST("/webhook/telegram", telegramBot.HandleWebhook)
		log.Println("Telegram Bot 已啟用")
	}

	// WebSocket 路由
	r.GET("/ws", handlers.HandleWebSocket)

	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// 啟動伺服器
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("伺服器啟動於 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}
