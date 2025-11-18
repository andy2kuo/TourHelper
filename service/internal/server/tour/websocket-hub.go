package tour

import (
	"sync"

	"github.com/andy2kuo/TourHelper/internal/logger"
)

// BroadcastMessage 廣播訊息結構
type BroadcastMessage struct {
	message []byte
	sender  *Client
	target  string // 目標客戶端 ID，空字串表示廣播給所有人
}

// Hub 管理所有 WebSocket 客戶端連線
type Hub struct {
	// 已註冊的客戶端
	clients map[*Client]bool

	// 客戶端 ID 對應表（用於點對點傳送）
	clientsByID map[string]*Client

	// 廣播訊息到所有客戶端
	broadcast chan BroadcastMessage

	// 註冊新客戶端
	register chan *Client

	// 取消註冊客戶端
	unregister chan *Client

	// 互斥鎖
	mu sync.RWMutex
}

// NewHub 建立新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		clientsByID: make(map[string]*Client),
		broadcast:   make(chan BroadcastMessage, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
	}
}

// Run 啟動 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if client.ID != "" {
				h.clientsByID[client.ID] = client
				logger.Infof("WebSocket 客戶端已註冊: %s (總共 %d 個連線)", client.ID, len(h.clients))
			} else {
				logger.Infof("WebSocket 客戶端已註冊 (總共 %d 個連線)", len(h.clients))
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.ID != "" {
					delete(h.clientsByID, client.ID)
					logger.Infof("WebSocket 客戶端已取消註冊: %s (剩餘 %d 個連線)", client.ID, len(h.clients))
				} else {
					logger.Infof("WebSocket 客戶端已取消註冊 (剩餘 %d 個連線)", len(h.clients))
				}
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			if message.target != "" {
				// 點對點傳送
				if targetClient, ok := h.clientsByID[message.target]; ok {
					select {
					case targetClient.send <- message.message:
					default:
						// 無法發送，關閉連線
						close(targetClient.send)
						delete(h.clients, targetClient)
						delete(h.clientsByID, targetClient.ID)
						logger.Warnf("無法發送訊息給客戶端 %s，連線已關閉", message.target)
					}
				} else {
					logger.Warnf("找不到目標客戶端: %s", message.target)
				}
			} else {
				// 廣播給所有客戶端（除了發送者）
				for client := range h.clients {
					if client == message.sender {
						continue // 不發送給自己
					}
					select {
					case client.send <- message.message:
					default:
						// 無法發送，關閉連線
						close(client.send)
						delete(h.clients, client)
						if client.ID != "" {
							delete(h.clientsByID, client.ID)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToAll 廣播訊息給所有客戶端
func (h *Hub) BroadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// 無法發送，關閉連線
			close(client.send)
			delete(h.clients, client)
			if client.ID != "" {
				delete(h.clientsByID, client.ID)
			}
		}
	}
}

// SendToClient 發送訊息給特定客戶端
func (h *Hub) SendToClient(clientID string, message []byte) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clientsByID[clientID]; ok {
		select {
		case client.send <- message:
			return nil
		default:
			return nil
		}
	}
	return nil
}

// GetClientCount 取得當前連線的客戶端數量
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClientIDs 取得所有客戶端 ID
func (h *Hub) GetClientIDs() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	ids := make([]string, 0, len(h.clientsByID))
	for id := range h.clientsByID {
		ids = append(ids, id)
	}
	return ids
}
