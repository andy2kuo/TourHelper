package tour

import (
	"encoding/json"
	"time"

	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/gorilla/websocket"
)

const (
	// 訊息傳送超時時間
	writeWait = 10 * time.Second

	// 接收 pong 訊息的超時時間
	pongWait = 60 * time.Second

	// ping 訊息發送間隔（必須小於 pongWait）
	pingPeriod = (pongWait * 9) / 10

	// 允許的最大訊息大小
	maxMessageSize = 512 * 1024 // 512 KB
)

// Client 代表一個 WebSocket 客戶端連線
type Client struct {
	// WebSocket 連線
	conn *websocket.Conn

	// 傳送訊息的通道
	send chan []byte

	// 所屬的 Hub
	hub *Hub

	// 客戶端 ID（可選，用於識別用戶）
	ID string

	// 客戶端元資料（可選，例如用戶資訊、房間資訊等）
	Metadata map[string]interface{}
}

// Message WebSocket 訊息格式
type Message struct {
	Type    string                 `json:"type"`              // 訊息類型
	Data    interface{}            `json:"data"`              // 訊息資料
	From    string                 `json:"from,omitempty"`    // 發送者 ID
	To      string                 `json:"to,omitempty"`      // 接收者 ID（空表示廣播）
	Payload map[string]interface{} `json:"payload,omitempty"` // 額外資料
}

// readPump 從 WebSocket 連線讀取訊息並傳送到 Hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf("WebSocket 讀取錯誤: %v", err)
			}
			break
		}

		// 解析訊息
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			logger.Warnf("無法解析 WebSocket 訊息: %v", err)
			continue
		}

		// 設定發送者
		msg.From = c.ID

		// 重新序列化訊息
		data, err := json.Marshal(msg)
		if err != nil {
			logger.Errorf("無法序列化訊息: %v", err)
			continue
		}

		// 傳送到 Hub 進行廣播或點對點傳送
		c.hub.broadcast <- BroadcastMessage{
			message: data,
			sender:  c,
			target:  msg.To,
		}
	}
}

// writePump 從通道讀取訊息並寫入 WebSocket 連線
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 關閉了通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 將排隊的訊息一起發送
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage 發送訊息給客戶端
func (c *Client) SendMessage(msgType string, data interface{}) error {
	msg := Message{
		Type: msgType,
		Data: data,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- msgBytes:
		return nil
	default:
		return nil
	}
}
