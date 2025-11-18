# Lobby Server 與 Tour Server 整合範例

本文件說明如何整合 Lobby Server 和 Tour Server,實現會員登入驗證與即時通訊的雙向溝通。

## 架構說明

```
[會員] --> [Lobby Server] <--> [Tour Server] <--> [WebSocket 客戶端]
            (登入驗證)           (即時通訊)
```

### Lobby Server
- **主要功能**: 處理會員登入驗證
- **端口**: 預設 8081
- **職責**:
  - 會員登入/登出
  - LINE 第三方登入驗證(待實作)
  - Token 驗證
  - 會員資訊管理
  - 與 Tour Server 雙向溝通

### Tour Server
- **主要功能**: 處理旅遊相關業務邏輯
- **端口**: 預設 8080
- **職責**:
  - HTTP + WebSocket 服務
  - 旅遊推薦
  - 景點查詢
  - 即時通訊(WebSocket)
  - Bot 整合(Line/Telegram)

## 啟動方式

### 方式 1: 分別啟動(開發模式)

```bash
# 終端機 1: 啟動 Tour Server
cd service
go run cmd/tour/main.go

# 終端機 2: 啟動 Lobby Server
cd service
go run cmd/lobby/main.go
```

### 方式 2: 整合啟動(生產模式)

建立一個整合的 main.go,同時啟動兩個伺服器:

```go
// cmd/integrated/main.go
package main

import (
    "github.com/andy2kuo/TourHelper/internal/config"
    "github.com/andy2kuo/TourHelper/internal/logger"
    "github.com/andy2kuo/TourHelper/internal/server/lobby"
    "github.com/andy2kuo/TourHelper/internal/server/tour"
)

func main() {
    // 初始化設定和日誌
    cfg, _ := config.Load("tourhelper", "production", "1.0.0")
    logger.Init("tourhelper", "production", cfg.Log)
    defer logger.GetLogger().Close()

    // 建立 Tour Server(埠號 8080)
    tourCfg := *cfg
    tourCfg.Server.Port = 8080
    tourServer := &tour.TourServer{}
    tourServer.Init(&server.Options{
        Config: &tourCfg,
        ServiceName: "tour_server",
        ServiceEnv: "production",
        Version: "1.0.0",
    })
    go tourServer.Start()

    // 建立 Lobby Server(埠號 8081)
    lobbyCfg := *cfg
    lobbyCfg.Server.Port = 8081
    lobbyServer := &lobby.LobbyServer{}
    lobbyServer.Init(&server.Options{
        Config: &lobbyCfg,
        ServiceName: "lobby_server",
        ServiceEnv: "production",
        Version: "1.0.0",
    })

    // 設定 Lobby 與 Tour 的連線
    lobbyServer.SetTourServer(tourServer)

    // 啟動 Lobby Server(會阻塞)
    lobbyServer.Start()
}
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

## Lobby 與 Tour 雙向溝通範例

### 從 Lobby 發送訊息到 Tour

```go
// 在 Lobby Server 中
func (s *LobbyServer) notifyMemberLogin(memberID string) {
    message, _ := json.Marshal(map[string]interface{}{
        "type": "member_login",
        "member_id": memberID,
        "timestamp": time.Now().Unix(),
    })

    // 廣播給所有 WebSocket 客戶端
    s.NotifyTourServer(message)

    // 或發送給特定會員
    s.SendToMember(memberID, message)
}
```

### 從 Tour 存取 Lobby 資訊

```go
// 在需要時,Tour Server 可以透過 HTTP API 呼叫 Lobby Server
resp, err := http.Post("http://localhost:8081/auth/verify",
    "application/json",
    bytes.NewBuffer(tokenJSON))
```

## 待實作功能清單

### Lobby Server

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

- [ ] 實作旅遊推薦 API
- [ ] 景點資訊查詢 API
- [ ] 使用者偏好設定 API
- [ ] WebSocket 訊息類型擴充

### 整合功能

- [ ] 會員登入後自動建立 WebSocket 連線
- [ ] 登出時自動關閉 WebSocket 連線
- [ ] 會員狀態同步(Lobby <-> Tour)
- [ ] 統一的錯誤處理機制

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

### Lobby 無法連線到 Tour

檢查:
1. Tour Server 是否已啟動
2. 埠號是否正確
3. 防火牆設定

### WebSocket 連線失敗

檢查:
1. 瀏覽器是否支援 WebSocket
2. URL 格式是否正確(`ws://` 或 `wss://`)
3. CORS 設定

### LINE 登入失敗(待實作時參考)

檢查:
1. LINE Channel ID 和 Secret 是否正確
2. Callback URL 是否已在 LINE Developer Console 設定
3. state 參數是否正確傳遞
