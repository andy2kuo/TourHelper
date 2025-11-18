package tour

import (
	"net/http"

	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允許所有來源（生產環境應該設定適當的檢查）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketHandler WebSocket 連線處理器
type WebSocketHandler struct {
	hub *Hub
}

// NewWebSocketHandler 建立 WebSocket 處理器
func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleWebSocket 處理 WebSocket 連線請求
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// 升級 HTTP 連線為 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Errorf("WebSocket 升級失敗: %v", err)
		return
	}

	// 從查詢參數或 Header 取得客戶端 ID（可選）
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = c.GetHeader("X-Client-ID")
	}

	// 建立新客戶端
	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      h.hub,
		ID:       clientID,
		Metadata: make(map[string]interface{}),
	}

	// 註冊客戶端
	h.hub.register <- client

	// 啟動讀寫協程
	go client.writePump()
	go client.readPump()

	logger.Infof("新的 WebSocket 連線建立: %s", clientID)
}

// HandleWebSocketInfo 提供 WebSocket 連線資訊的 HTTP API
func (h *WebSocketHandler) HandleWebSocketInfo(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":       "ok",
		"clients":      h.hub.GetClientCount(),
		"client_ids":   h.hub.GetClientIDs(),
		"endpoint":     "/ws",
		"description":  "WebSocket endpoint for real-time communication",
		"query_params": "client_id (optional) - Unique identifier for the client",
	})
}
