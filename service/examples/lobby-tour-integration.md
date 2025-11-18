# Lobby Server 與 Tour Server 整合範例

本文件說明如何整合 Lobby Server 和 Tour Server，使用 Redis 作為非同步狀態交換中心。

## 架構說明

```
[會員] --> [Lobby Server] --> [Redis] <-- [Tour Server] <--> [WebSocket 客戶端]
            (登入驗證)       (狀態中心)    (即時通訊)
```

**關鍵特性**：
- Lobby 和 Tour 不直接通訊
- Tour Server 更新狀態到 Redis
- Lobby Server 定時從 Redis 讀取 Tour 狀態
- 使用 Redis Pub/Sub 或輪詢機制實現非同步狀態同步

### Lobby Server
- **主要功能**: 處理會員登入驗證和會員管理
- **端口**: 預設 8081
- **職責**:
  - 會員登入/登出
  - LINE 第三方登入驗證(待實作)
  - Token 驗證
  - 會員資訊管理
  - 將會員狀態更新到 Redis
  - 定時從 Redis 取得 Tour Server 狀態

### Tour Server
- **主要功能**: 處理旅遊相關業務邏輯和即時通訊
- **端口**: 預設 8080
- **職責**:
  - HTTP + WebSocket 服務
  - 旅遊推薦
  - 景點查詢
  - 即時通訊(WebSocket)
  - Bot 整合(Line/Telegram)
  - 將伺服器狀態和旅遊資訊更新到 Redis

### Backend Admin Server
- **主要功能**: 後台管理功能
- **端口**: 預設 8082
- **職責**:
  - 管理員登入驗證
  - 會員管理(查詢、狀態、刪除)
  - Tour Server 管理
  - 系統設定
  - 日誌查詢

## 啟動方式

### Windows 開發模式（推薦）

使用提供的批次檔啟動各個伺服器：

```batch
# 啟動 Tour Server (預設 8080 埠)
cd service
runTour.bat

# 啟動 Lobby Server (預設 8081 埠)
cd service
runLobby.bat

# 啟動 Backend Admin Server (預設 8082 埠)
cd service
runBackend.bat
```

每個批次檔會提示輸入：
- **Service Name**: 服務名稱（預設值已提供）
- **Service Version**: 版本號（預設 0.0.1-dev）
- **Service Environment**: 環境（dev/staging/production）

### Linux/Mac 開發模式

```bash
# 終端機 1: 啟動 Tour Server
cd service
go run -ldflags "-X main.SERVICE_NAME=tour_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/tour/main.go

# 終端機 2: 啟動 Lobby Server
cd service
go run -ldflags "-X main.SERVICE_NAME=lobby_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/lobby/main.go

# 終端機 3: 啟動 Backend Admin Server
cd service
go run -ldflags "-X main.SERVICE_NAME=backend_admin_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/backend_admin/main.go
```

### Redis 啟動

三個伺服器都需要 Redis 支援（待實作整合）：

```bash
# Docker 方式
docker run -d -p 6379:6379 redis:latest

# 或使用 Docker Compose
docker-compose up -d redis
```

## API 使用範例

### Lobby Server API

#### 1. 會員登入

```bash
curl -X POST http://localhost:8081/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "platform": "web"
  }'
```

**回應**:
```json
{
  "success": true,
  "message": "登入成功(功能待實作)",
  "token": "dummy_token",
  "member_id": "testuser"
}
```

#### 2. 驗證 Token

```bash
curl -X POST http://localhost:8081/auth/verify \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your_jwt_token"
  }'
```

#### 3. 取得會員資訊

```bash
curl http://localhost:8081/member/user123
```

#### 4. 更新會員資訊

```bash
curl -X PUT http://localhost:8081/member/user123 \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "新暱稱",
    "avatar": "https://example.com/avatar.jpg"
  }'
```

#### 5. 登出

```bash
curl -X POST http://localhost:8081/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your_jwt_token"
  }'
```

### Tour Server API

#### 1. 健康檢查

```bash
curl http://localhost:8080/health
```

#### 2. WebSocket 連線資訊

```bash
curl http://localhost:8080/ws/info
```

#### 3. WebSocket 連線

```javascript
// 在瀏覽器或 Node.js 中
const ws = new WebSocket('ws://localhost:8080/ws?client_id=user123');

ws.onopen = () => {
  console.log('WebSocket 連線成功');

  // 發送訊息
  ws.send(JSON.stringify({
    type: 'recommendation_request',
    data: {
      latitude: 25.0330,
      longitude: 121.5654,
      max_distance: 50
    }
  }));
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('收到訊息:', message);
};
```

## Redis 狀態交換範例

### Lobby Server 更新會員狀態到 Redis

```go
// 在 Lobby Server 中（待實作）
func (s *LobbyServer) UpdateMemberStatusToRedis(memberID string, status map[string]interface{}) error {
    // TODO: 實作 Redis 狀態更新
    // 1. 連接 Redis
    // 2. 設定 key: member:{memberID}:status
    // 3. 儲存狀態資料(JSON 格式)
    // 4. 設定過期時間(例如: 30 分鐘)

    // 範例實作:
    // key := fmt.Sprintf("member:%s:status", memberID)
    // data, _ := json.Marshal(status)
    // return s.redisClient.Set(ctx, key, data, 30*time.Minute).Err()

    return nil
}

// 會員登入時更新狀態
func (s *LobbyServer) handleLogin(c *gin.Context) {
    // ... 驗證邏輯 ...

    // 登入成功後更新到 Redis
    s.UpdateMemberStatusToRedis(memberID, map[string]interface{}{
        "status": "online",
        "login_time": time.Now().Unix(),
        "platform": platform,
    })
}
```

### Lobby Server 定時取得 Tour 狀態

```go
// 在 Lobby Server 中（待實作）
func (s *LobbyServer) GetTourStatusFromRedis() (map[string]interface{}, error) {
    // TODO: 實作從 Redis 讀取 Tour Server 狀態
    // 1. 連接 Redis
    // 2. 讀取 key: tour:server:status
    // 3. 解析 JSON 資料

    // 範例實作:
    // val, err := s.redisClient.Get(ctx, "tour:server:status").Result()
    // if err != nil {
    //     return nil, err
    // }
    // var status map[string]interface{}
    // json.Unmarshal([]byte(val), &status)
    // return status, nil

    return map[string]interface{}{
        "status": "unknown",
        "message": "Redis integration not implemented",
    }, nil
}

// 定時監控 Tour Server 狀態
func (s *LobbyServer) StartTourStatusMonitor() {
    // TODO: 實作定時監控
    // ticker := time.NewTicker(10 * time.Second)
    // go func() {
    //     for range ticker.C {
    //         status, err := s.GetTourStatusFromRedis()
    //         if err != nil {
    //             logger.Errorf("取得 Tour 狀態失敗: %v", err)
    //             continue
    //         }
    //         logger.Infof("Tour Server 狀態: %+v", status)
    //     }
    // }()
}
```

### Tour Server 更新狀態到 Redis

```go
// 在 Tour Server 中（待實作）
func (s *TourServer) UpdateStatusToRedis() error {
    // TODO: 實作 Redis 狀態更新
    // 1. 收集伺服器狀態資訊
    // 2. 設定 key: tour:server:status
    // 3. 儲存狀態資料

    // 範例實作:
    // status := map[string]interface{}{
    //     "status": "running",
    //     "ws_clients": s.wsHub.GetClientCount(),
    //     "update_time": time.Now().Unix(),
    // }
    // data, _ := json.Marshal(status)
    // return s.redisClient.Set(ctx, "tour:server:status", data, 1*time.Minute).Err()

    return nil
}

// 定時更新狀態
func (s *TourServer) StartStatusReporter() {
    // TODO: 每 5 秒更新一次狀態到 Redis
}
```

## 待實作功能清單

### Lobby Server

- [ ] **Redis 客戶端整合**
  - [ ] 建立 Redis 連線池
  - [ ] 實作 `UpdateMemberStatusToRedis()` 方法
  - [ ] 實作 `GetTourStatusFromRedis()` 方法
  - [ ] 實作 `StartTourStatusMonitor()` 定時監控
- [ ] 實作真正的帳號密碼驗證
- [ ] 實作 JWT Token 產生與驗證
- [ ] **實作 LINE 第三方登入驗證** (已預留架構)
  - 參考: https://developers.line.biz/en/docs/line-login/
  - OAuth 2.0 流程
  - Access Token 換取
  - 使用者資訊取得
- [ ] 資料庫整合(會員資料儲存)
- [ ] Session 管理
- [ ] 會員線上狀態追蹤

### Tour Server

- [ ] **Redis 客戶端整合**
  - [ ] 建立 Redis 連線池
  - [ ] 實作 `UpdateStatusToRedis()` 方法
  - [ ] 實作 `StartStatusReporter()` 定時上報狀態
  - [ ] 實作從 Redis 讀取會員狀態
- [ ] 實作旅遊推薦 API
- [ ] 景點資訊查詢 API
- [ ] 使用者偏好設定 API
- [ ] WebSocket 訊息類型擴充

### Backend Admin Server

- [ ] **Redis 客戶端整合**
  - [ ] 建立 Redis 連線池
  - [ ] 實作從 Redis 讀取伺服器狀態
  - [ ] 實作 Redis 快取管理功能
- [ ] 實作管理員登入驗證
- [ ] 實作管理員權限管理(角色、權限)
- [ ] 實作會員管理 CRUD
- [ ] 實作 Tour Server 管理功能
- [ ] 實作系統設定管理
- [ ] 實作日誌查詢功能

### Redis 架構設計

- [ ] 設計 Redis Key 命名規範
  - `member:{memberID}:status` - 會員狀態
  - `tour:server:status` - Tour Server 狀態
  - `tour:destinations` - 景點快取
  - `session:{token}` - Session 資料
- [ ] 設計資料過期策略
- [ ] 實作 Redis Pub/Sub (可選)
  - 即時通知會員狀態變更
  - 即時通知伺服器狀態變更

### 整合功能

- [ ] 會員登入後自動建立 WebSocket 連線
- [ ] 登出時自動關閉 WebSocket 連線
- [ ] 會員狀態透過 Redis 同步
- [ ] 統一的錯誤處理機制
- [ ] 統一的日誌格式

## 安全性建議

1. **HTTPS**: 生產環境必須使用 HTTPS
2. **CORS**: 適當設定 CORS 政策
3. **Rate Limiting**: 實作請求頻率限制
4. **Input Validation**: 嚴格驗證所有輸入資料
5. **Token 安全**:
   - 使用強密碼簽署 JWT
   - 設定適當的過期時間
   - 實作 Token 刷新機制
6. **LINE 登入**:
   - 驗證 state 參數防止 CSRF
   - 使用 HTTPS callback URL
   - 驗證 ID Token

## 測試

使用提供的測試工具:

```bash
# WebSocket 測試
firefox service/examples/websocket-test.html

# API 測試
curl http://localhost:8081/health
curl http://localhost:8080/health
```

## 疑難排解

### Redis 連線失敗

檢查:
1. Redis 伺服器是否已啟動（`docker ps` 或 `redis-cli ping`）
2. Redis 連線參數是否正確（host、port、password）
3. 防火牆設定是否允許連線
4. 網路連線是否正常

### Lobby 無法取得 Tour 狀態

檢查:
1. Tour Server 是否已啟動
2. Tour Server 是否正確將狀態寫入 Redis
3. Redis Key 是否正確（`tour:server:status`）
4. 使用 `redis-cli` 檢查資料：`GET tour:server:status`

### WebSocket 連線失敗

檢查:
1. 瀏覽器是否支援 WebSocket
2. URL 格式是否正確(`ws://` 或 `wss://`)
3. CORS 設定
4. Tour Server 是否正確啟動

### LINE 登入失敗(待實作時參考)

檢查:
1. LINE Channel ID 和 Secret 是否正確
2. Callback URL 是否已在 LINE Developer Console 設定
3. state 參數是否正確傳遞

### 伺服器啟動腳本執行錯誤

檢查:
1. 是否在 `service/` 目錄下執行
2. `go.mod` 檔案是否存在
3. Go 環境是否正確設定
4. 依賴套件是否已安裝（`go mod download`）

## Redis 管理工具

### 使用 redis-cli 檢查狀態

```bash
# 連接到 Redis
redis-cli

# 查看所有 keys
KEYS *

# 查看 Tour Server 狀態
GET tour:server:status

# 查看會員狀態
GET member:user123:status

# 查看 key 過期時間
TTL tour:server:status

# 刪除特定 key
DEL tour:server:status
```

### 使用 Redis Commander (Web UI)

```bash
# 使用 Docker 啟動
docker run -d -p 8083:8081 --name redis-commander --env REDIS_HOSTS=local:host.docker.internal:6379 rediscommander/redis-commander

# 瀏覽器開啟
# http://localhost:8083
```
