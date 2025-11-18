# TourHelper Backend

TourHelper 後端服務，使用 Go 語言和 Gin 框架建構，提供旅遊推薦 API、Bot 整合等功能。

## 技術架構

- **程式語言**：Go 1.25+
- **Web 框架**：Gin Web Framework
- **WebSocket**：gorilla/websocket
- **資料庫**：支援 SQLite / MySQL / PostgreSQL
- **ORM**：GORM
- **配置管理**：Viper
- **Bot 整合**：Line Bot SDK、Telegram Bot API
- **測試框架**：Go testing、Testify
- **日誌**：Logrus + Lumberjack（支援 log rotation）

## 專案結構

```text
service/
├── cmd/                    # 主程式進入點
│   ├── tour/               # Tour Server (旅遊推薦服務)
│   │   └── main.go         # Tour 伺服器啟動檔案，初始化設定、建立 tour server、優雅關閉
│   ├── lobby/              # Lobby Server (會員登入驗證服務)
│   │   └── main.go         # Lobby 伺服器啟動檔案，初始化設定、建立 lobby server、優雅關閉
│   └── backend/      # Backend Server (後台管理服務)
│       └── main.go         # Backend Admin 伺服器啟動檔案，初始化設定、建立 backend admin server、優雅關閉
├── internal/               # 私有應用程式碼
│   ├── config/             # 設定管理
│   │   └── config.go       # 使用 Viper 管理設定，支援 YAML 和環境變數
│   ├── server/             # 伺服器實作
│   │   ├── server.go       # Server 介面定義
│   │   ├── tour/           # Tour Server 實作（HTTP + WebSocket）
│   │   │   ├── tour-server.go       # Tour 伺服器實作（Gin）
│   │   │   ├── tour-handler.go      # Tour 請求處理器
│   │   │   ├── websocket-hub.go     # WebSocket 連線管理（Hub）
│   │   │   ├── websocket-client.go  # WebSocket 客戶端連線
│   │   │   └── websocket-handler.go # WebSocket 請求處理器
│   │   ├── lobby/          # Lobby Server 實作（HTTP Only）
│   │   │   ├── lobby-server.go      # Lobby 伺服器實作（Gin）
│   │   │   └── lobby-handler.go     # Lobby 請求處理器（登入、會員管理）
│   │   └── backend/  # Backend Server 實作（HTTP Only）
│   │       ├── backend-admin-server.go   # Backend Admin 伺服器實作（Gin）
│   │       └── backend-admin-handler.go  # Backend Admin 請求處理器（管理功能）
│   ├── database/           # 資料庫管理
│   │   ├── database.go     # 資料庫連線管理、初始化
│   │   └── example_usage.go # 使用範例
│   ├── logger/             # 日誌管理
│   │   ├── logger.go       # Logrus + Lumberjack，支援 log rotation
│   │   ├── umask_unix.go   # Unix/Linux 平台的檔案權限設定
│   │   └── umask_windows.go # Windows 平台的檔案權限設定
│   ├── services/           # 業務邏輯服務層
│   │   ├── services.go     # Services 單例管理
│   │   ├── recommendation_service.go  # 推薦服務
│   │   └── weather_service.go         # 天氣服務
│   ├── dao/                # 資料庫存取層（Data Access Object）
│   │   ├── dao.go          # DAO 單例管理
│   │   ├── user_dao.go     # 使用者 CRUD 操作
│   │   └── destination_dao.go  # 景點 CRUD 操作
│   ├── models/             # 資料模型
│   │   ├── models.go       # 定義 User, Destination, Tag 等資料結構
│   │   └── database.go     # 資料庫初始化、遷移和範例資料填充
│   └── bot/                # Bot 整合
│       ├── line/           # Line Bot 實作
│       │   └── line.go     # 處理 Line webhook、訊息回覆
│       └── telegram/       # Telegram Bot 實作
│           └── telegram.go # 處理 Telegram webhook、指令處理
├── pkg/                    # 可重用的公開函式庫
│   └── utils/              # 工具函式
│       ├── utils.go        # 距離計算、時間估算等工具
│       └── utils_test.go   # 工具函式測試
├── configs/                # 設定檔
│   └── config.example.yaml # 設定檔範例
├── go.mod                  # Go 模組依賴
├── go.sum                  # Go 模組校驗
├── .gitignore              # Git 忽略檔案
└── README.md               # 本檔案
```

## 三伺服器架構

本專案採用三個獨立伺服器架構，使用 Redis 進行狀態同步：

```text
                    [Redis 狀態中心]
                           ↑  ↓
         ┌─────────────────┴──┴───────────────┐
         ↓                  ↓                  ↓
[Tour Server]      [Lobby Server]    [Backend Server]
(HTTP+WebSocket)      (HTTP Only)         (HTTP Only)
  8080                  8081                 8082

  - 旅遊推薦          - 會員登入驗證       - 後台管理
  - 景點查詢          - LINE 登入          - 會員管理
  - 即時通訊          - Token 驗證          - 系統設定
  - Bot 整合          - 會員資訊管理        - 日誌查詢
```

### 通訊方式

- **Tour ↔ Lobby**: 透過 Redis 非同步通訊（不直接呼叫）
  - Tour 將狀態寫入 Redis
  - Lobby 定時從 Redis 讀取 Tour 狀態
- **Backend Admin ↔ Redis**: 讀取各伺服器狀態，管理系統設定

## 分層架構

本專案採用經典的分層架構設計，確保關注點分離和程式碼可維護性：

```text
外部請求（HTTP / WebSocket / Bot 等）
    ↓
[Handlers] ← server/tour/tour-handler.go
           ← server/lobby/lobby-handler.go
           ← server/backend/backend-admin-handler.go
    ↓ 處理請求/回應、參數驗證
[Services] ← services/
    ↓ 業務邏輯處理
[DAO] ← dao/
    ↓ 資料庫 CRUD 操作
[Models] ← models/
    ↓ 資料結構定義
資料庫 / Redis
```

### 各層職責

1. **Handlers（處理器層）**
   - 位置：
     - `internal/server/tour/tour-handler.go` - Tour Server 請求處理
     - `internal/server/lobby/lobby-handler.go` - Lobby Server 請求處理（登入、會員管理）
     - `internal/server/backend/backend-admin-handler.go` - Backend Admin 請求處理
   - 職責：處理各種請求和回應、參數解析和驗證、呼叫 Service 層
   - 範例：登入驗證、旅遊推薦、WebSocket 升級、Bot webhook

2. **Services（業務邏輯層）**
   - 位置：`internal/services/`
   - 職責：實作核心業務邏輯、協調多個 DAO、處理複雜的業務規則
   - 特點：使用單例模式（`services.Get()`）
   - 範例：RecommendationService、WeatherService

3. **DAO（資料存取層）**
   - 位置：`internal/dao/`
   - 職責：封裝所有資料庫 CRUD 操作、處理查詢條件
   - 特點：使用單例模式（`dao.Get()`）
   - 範例：UserDAO、DestinationDAO

4. **Models（資料模型層）**
   - 位置：`internal/models/`
   - 職責：定義資料結構、GORM 標籤、資料庫遷移
   - 範例：User、Destination、Tag

### 依賴關係

- Handlers → Services → DAO → Models
- 每一層只能依賴下一層，不能跨層調用
- 使用介面定義行為，便於測試和替換實作

## 核心模組說明

### cmd/tour/main.go

Tour Server 進入點，負責：

- 載入設定（使用 config.Load）
- 初始化 Logger
- 建立 Tour Server 實例
- 啟動 Tour Server（HTTP + WebSocket）
- 處理優雅關閉（Graceful Shutdown）

### cmd/lobby/main.go

Lobby Server 進入點，負責：

- 載入設定（使用 config.Load）
- 初始化 Logger
- 建立 Lobby Server 實例
- 啟動 Lobby Server（HTTP Only）
- 處理優雅關閉（Graceful Shutdown）

### cmd/backend/main.go

Backend Server 進入點，負責：

- 載入設定（使用 config.Load）
- 初始化 Logger
- 建立 Backend Server 實例
- 啟動 Backend Server（HTTP Only）
- 處理優雅關閉（Graceful Shutdown）

### internal/server/

伺服器實作模組，支援多種伺服器類型：

- **server.go**：定義 Server 介面，規範所有伺服器類型的行為
  - `Init(opts *Options)`：初始化伺服器
  - `Start()`：啟動伺服器
  - `Stop()`：停止伺服器
  - `Name()`：返回伺服器名稱

- **tour/**：Tour Server 實作（HTTP + WebSocket）
  - **tour-server.go**：使用 Gin 框架實作 Tour HTTP/WebSocket 伺服器
    - 註冊旅遊相關 API 路由
    - 整合 WebSocket Hub
    - 整合 Line 和 Telegram Bot webhook
    - 支援優雅關閉
  - **tour-handler.go**：Tour 請求處理器
    - HealthCheckHandler：處理健康檢查請求
  - **websocket-hub.go**：WebSocket 連線管理中心
    - 管理所有 WebSocket 客戶端連線
    - 支援廣播訊息給所有客戶端
    - 支援點對點訊息傳送
    - 自動清理斷線的客戶端
  - **websocket-client.go**：WebSocket 客戶端連線
    - 處理單一客戶端的讀寫操作
    - 支援心跳檢測（ping/pong）
    - 自動處理訊息序列化
  - **websocket-handler.go**：WebSocket HTTP 請求處理器
    - 處理 WebSocket 升級請求
    - 提供連線資訊 API

- **lobby/**：Lobby Server 實作（HTTP Only）
  - **lobby-server.go**：使用 Gin 框架實作 Lobby HTTP 伺服器
    - 註冊會員驗證相關路由
    - Redis 狀態管理（待實作）
    - 支援優雅關閉
  - **lobby-handler.go**：Lobby 請求處理器
    - handleLogin：一般登入驗證
    - handleLineLogin：LINE 第三方登入（待實作）
    - handleLogout：登出
    - handleVerifyToken：Token 驗證
    - handleGetMemberInfo：取得會員資訊
    - handleUpdateMemberInfo：更新會員資訊

- **backend/**：Backend Server 實作（HTTP Only）
  - **backend-admin-server.go**：使用 Gin 框架實作 Backend Admin HTTP 伺服器
    - 註冊管理相關路由
    - Redis 狀態讀取（待實作）
    - 支援優雅關閉
  - **backend-admin-handler.go**：Backend Admin 請求處理器
    - handleAdminLogin：管理員登入
    - handleGetMembers：取得會員列表
    - handleGetMemberDetail：取得會員詳情
    - handleUpdateMemberStatus：更新會員狀態
    - handleDeleteMember：刪除會員
    - handleGetTourStatus：取得 Tour Server 狀態
    - handleGetDestinations：取得景點列表
    - 其他管理功能處理器

**設計理念**：透過 Server 介面和三伺服器架構，可以：

1. 獨立部署和擴展各個伺服器
2. 使用 Redis 進行非同步狀態同步，降低耦合度
3. 清晰的職責劃分：
   - Tour Server：旅遊推薦和即時通訊
   - Lobby Server：會員驗證和管理
   - Backend Server：後台管理功能
4. 輕鬆新增其他類型的伺服器（例如：gRPC Server 等）

### internal/config/config.go

設定管理模組，使用 Viper 載入：

- YAML 設定檔（`configs/config.yaml`）
- 環境變數
- 支援多環境配置

### internal/models/

資料模型層，包含：

- `models.go`：定義所有資料結構（User, Destination, Tag, Preference 等），包含 GORM 標籤
- `database.go`：資料庫連線初始化、自動遷移、範例資料填充功能

### internal/services/

業務邏輯層，負責處理核心業務邏輯：

- **services.go**：Services 單例管理
  - 使用 `services.Get()` 取得 services 實例
  - 集中管理所有 service 的初始化

- **recommendation_service.go**：推薦服務
  - 實作推薦演算法核心邏輯
  - 根據位置、天氣、距離計算適合度評分
  - 協調多個 DAO 取得資料

- **weather_service.go**：天氣服務
  - 整合第三方天氣 API
  - 處理天氣資料的業務邏輯

### internal/dao/

資料存取層（Data Access Object），負責所有資料庫 CRUD 操作：

- **dao.go**：DAO 單例管理
  - 使用 `dao.Get()` 取得 DAO 實例
  - 需先呼叫 `dao.SetDB(db)` 設定資料庫連線
  - 集中管理所有 DAO 的初始化

- **user_dao.go**：使用者 DAO
  - 封裝使用者資料的 CRUD 操作
  - 提供查詢、建立、更新、刪除等方法

- **destination_dao.go**：景點 DAO
  - 封裝景點資料的 CRUD 操作
  - 提供地理位置相關查詢等方法

**使用範例**：

```go
// 初始化
dao.SetDB(db)

// 使用
daos := dao.Get()
services := services.Get()

// 呼叫 service
result := services.Recommendation.GetRecommendations(...)
```

### internal/bot/

Bot 整合層：

- `line/line.go`：處理 Line webhook、訊息解析、回覆訊息
- `telegram/telegram.go`：處理 Telegram webhook、指令處理、訊息回覆

### internal/logger/

日誌管理模組：

- 使用 Logrus 作為日誌框架
- 使用 Lumberjack 實現自動 log rotation
- 支援多級別日誌（Debug、Info、Warn、Error、Fatal）
- 同時輸出到檔案和 console
- 根據環境自動選擇日誌路徑：
  - **開發環境 (dev)**：`./log/{SERVICE_NAME}.log`
  - **生產環境 (Linux)**：`/var/log/{SERVICE_NAME}/{SERVICE_NAME}.log`
  - **生產環境 (Windows)**：`C:/ProgramData/{SERVICE_NAME}/log/{SERVICE_NAME}.log`
- 可設定檔案大小、保留數量和天數

### pkg/utils/

公開工具函式庫：

- 距離計算（Haversine 公式）
- 旅行時間估算
- 其他通用工具函式

## 開發指令

### 安裝依賴

```bash
# 下載依賴
go mod download

# 整理依賴
go mod tidy
```

### 執行程式

#### Windows 開發模式（推薦）

使用提供的批次檔啟動各個伺服器：

```batch
# 啟動 Tour Server (預設 8080 埠)
runTour.bat

# 啟動 Lobby Server (預設 8081 埠)
runLobby.bat

# 啟動 Backend Server (預設 8082 埠)
runBackend.bat
```

每個批次檔會提示輸入服務名稱、版本和環境。

#### Linux/Mac 開發模式

```bash
# 執行 Tour Server
go run -ldflags "-X main.SERVICE_NAME=tour_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/tour/main.go

# 執行 Lobby Server
go run -ldflags "-X main.SERVICE_NAME=lobby_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/lobby/main.go

# 執行 Backend Server
go run -ldflags "-X main.SERVICE_NAME=backend_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/backend/main.go
```

### 建置

```bash
# 基本建置
go build -o tour cmd/tour/main.go
go build -o lobby cmd/lobby/main.go
go build -o backend cmd/backend/main.go

# 建置到特定目錄
go build -o bin/tour cmd/tour/main.go
go build -o bin/lobby cmd/lobby/main.go
go build -o bin/backend cmd/backend/main.go

# 建置並注入版本資訊（使用 -ldflags）
# Tour Server
go build -ldflags "\
  -X main.SERVICE_NAME=tour_server \
  -X main.SERVICE_ENV=production \
  -X main.SERVICE_VERSION=1.0.0" \
  -o tour cmd/tour/main.go

# Lobby Server
go build -ldflags "\
  -X main.SERVICE_NAME=lobby_server \
  -X main.SERVICE_ENV=production \
  -X main.SERVICE_VERSION=1.0.0" \
  -o lobby cmd/lobby/main.go

# Backend Server
go build -ldflags "\
  -X main.SERVICE_NAME=backend_server \
  -X main.SERVICE_ENV=production \
  -X main.SERVICE_VERSION=1.0.0" \
  -o backend cmd/backend/main.go

# 開發環境建置
go build -ldflags "\
  -X main.SERVICE_NAME=tour_server \
  -X main.SERVICE_ENV=dev \
  -X main.SERVICE_VERSION=0.0.1-dev" \
  -o tour cmd/tour/main.go

go build -ldflags "\
  -X main.SERVICE_NAME=lobby_server \
  -X main.SERVICE_ENV=dev \
  -X main.SERVICE_VERSION=0.0.1-dev" \
  -o lobby cmd/lobby/main.go

go build -ldflags "\
  -X main.SERVICE_NAME=backend_server \
  -X main.SERVICE_ENV=dev \
  -X main.SERVICE_VERSION=0.0.1-dev" \
  -o backend cmd/backend/main.go

# 執行建置的檔案
./tour           # 啟動 Tour Server
./lobby          # 啟動 Lobby Server
./backend  # 啟動 Backend Server
```

**說明**：

- `-ldflags` 用於在編譯時注入變數值
- `SERVICE_NAME`：服務名稱（用於識別服務）
- `SERVICE_ENV`：運行環境（dev/staging/production）
- `SERVICE_VERSION`：版本號碼（建議使用語義化版本）

### 執行測試

```bash
# 執行所有測試
go test ./...

# 執行測試並顯示詳細輸出
go test -v ./...

# 執行測試並顯示覆蓋率
go test -cover ./...

# 產生覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 執行特定套件的測試
go test ./pkg/utils

# 執行特定測試函式
go test -run TestCalculateDistance ./pkg/utils
```

### 程式碼品質

```bash
# 格式化程式碼
go fmt ./...

# 執行 go vet（檢查常見錯誤）
go vet ./...

# 執行 linter（需先安裝 golangci-lint）
golangci-lint run

# 安裝 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 開發工具

#### Air - 自動重載

使用 [Air](https://github.com/cosmtrek/air) 進行自動重載：

```bash
# 安裝 Air
go install github.com/cosmtrek/air@latest

# 執行 Air（使用根目錄的 .air.toml 配置）
air
```

## 設定說明

### 設定檔位置

- 範例設定檔：`configs/config.example.yaml`
- 實際設定檔：`configs/config.yaml`（需自行建立）

### 設定結構

```yaml
# 伺服器設定
server:
  host: 0.0.0.0
  port: 8080
  mode: debug  # debug, release, test

# 資料庫設定（MySQL）
database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: tourhelper
  charset: utf8mb4
  parsetime: true
  loc: Local
  maxidleconns: 10      # 最大閒置連線數
  maxopenconns: 100     # 最大開啟連線數
  connmaxlifetime: 3600 # 連線最大生命週期（秒）

# Line Bot 設定
line:
  enabled: true
  channel_secret: YOUR_CHANNEL_SECRET
  channel_access_token: YOUR_CHANNEL_ACCESS_TOKEN

# Telegram Bot 設定
telegram:
  enabled: true
  token: YOUR_BOT_TOKEN

# 天氣 API 設定
weather:
  api_key: YOUR_API_KEY
  provider: openweathermap  # openweathermap, weatherapi

# 地圖 API 設定
maps:
  api_key: YOUR_API_KEY
  provider: google  # google, here, mapbox

# 日誌設定
log:
  maxsize: 100      # 單一日誌檔案最大大小（MB）
  maxbackups: 3     # 保留的舊日誌檔案數量
  maxage: 28        # 保留的天數
  compress: true    # 是否壓縮舊日誌檔案

# 推薦演算法參數
recommendation:
  default_max_distance: 50  # 預設最大距離（公里）
  weather_weight: 0.3       # 天氣因素權重
  distance_weight: 0.4      # 距離因素權重
  preference_weight: 0.3    # 偏好因素權重
```

### 環境變數

可透過環境變數覆蓋設定檔（使用 `TOURHELPER_` 前綴）：

```bash
# 資料庫設定
export TOURHELPER_DATABASE_HOST=localhost
export TOURHELPER_DATABASE_PORT=3306
export TOURHELPER_DATABASE_USER=root
export TOURHELPER_DATABASE_PASSWORD=your_password
export TOURHELPER_DATABASE_DBNAME=tourhelper

# Line Bot
export TOURHELPER_LINE_ENABLED=true
export TOURHELPER_LINE_CHANNEL_SECRET=your_secret
export TOURHELPER_LINE_CHANNEL_ACCESS_TOKEN=your_token

# Telegram Bot
export TOURHELPER_TELEGRAM_ENABLED=true
export TOURHELPER_TELEGRAM_TOKEN=your_token

# 天氣 API
export TOURHELPER_WEATHER_API_KEY=your_key
export TOURHELPER_WEATHER_PROVIDER=openweathermap
```

## API 端點

### 健康檢查

```http
GET /health
```

回應：

```json
{
  "status": "ok",
  "service": "TourHelper",
  "env": "dev",
  "version": "0.0.1"
}
```

### WebSocket

#### WebSocket 連線

```http
GET /ws?client_id={客戶端ID}
```

升級 HTTP 連線為 WebSocket，用於即時通訊。

**查詢參數**:

- `client_id`（可選）：客戶端唯一識別碼，用於點對點訊息傳送

**訊息格式**:

```json
{
  "type": "訊息類型",
  "data": "訊息內容",
  "from": "發送者ID",
  "to": "接收者ID（選填，空表示廣播）",
  "payload": {
    "額外資料": "值"
  }
}
```

**範例**：

```javascript
// 建立 WebSocket 連線
const ws = new WebSocket('ws://localhost:8080/ws?client_id=user123');

// 監聽連線開啟
ws.onopen = () => {
  console.log('WebSocket 連線成功');

  // 發送訊息
  ws.send(JSON.stringify({
    type: 'chat',
    data: { text: 'Hello World' }
  }));
};

// 監聽收到的訊息
ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('收到訊息:', message);
};
```

#### WebSocket 連線資訊

```http
GET /ws/info
```

取得當前 WebSocket 連線資訊。

回應：

```json
{
  "status": "ok",
  "clients": 5,
  "client_ids": ["user123", "user456"],
  "endpoint": "/ws",
  "description": "WebSocket endpoint for real-time communication",
  "query_params": "client_id (optional) - Unique identifier for the client"
}
```

### Webhook

#### Line Bot Webhook

```http
POST /webhook/line
```

處理 Line Messaging API 的 webhook 事件。

#### Telegram Bot Webhook

```http
POST /webhook/telegram
```

處理 Telegram Bot API 的 webhook 事件。

## 資料庫設計

### 連線池設定

MySQL 連線池參數說明：

- **MaxIdleConns**：最大閒置連線數（預設：10）
  - 保持在池中的閒置連線數量
  - 設定過低可能導致頻繁建立新連線
  - 設定過高會佔用過多資源

- **MaxOpenConns**：最大開啟連線數（預設：100）
  - 同時開啟的最大連線數
  - 建議根據資料庫伺服器能力設定
  - 避免設定過高導致資料庫過載

- **ConnMaxLifetime**：連線最大生命週期（預設：3600 秒）
  - 連線可以被重複使用的最長時間
  - 超過此時間的連線會被關閉並重新建立
  - 建議設定為略小於資料庫的 wait_timeout

### User 使用者

| 欄位 | 型別 | 說明 |
|------|------|------|
| id | uint | 主鍵 |
| external_id | string | 外部 ID（Line/Telegram） |
| platform | string | 平台（line/telegram/web） |
| created_at | time | 建立時間 |
| updated_at | time | 更新時間 |

### Destination 景點

| 欄位 | 型別 | 說明 |
|------|------|------|
| id | uint | 主鍵 |
| name | string | 景點名稱 |
| description | string | 景點描述 |
| latitude | float64 | 緯度 |
| longitude | float64 | 經度 |
| tags | []Tag | 標籤（多對多） |
| created_at | time | 建立時間 |
| updated_at | time | 更新時間 |

### Tag 標籤

| 欄位 | 型別 | 說明 |
|------|------|------|
| id | uint | 主鍵 |
| name | string | 標籤名稱 |
| created_at | time | 建立時間 |

### UserPreference 使用者偏好

| 欄位 | 型別 | 說明 |
|------|------|------|
| id | uint | 主鍵 |
| user_id | uint | 使用者 ID（外鍵） |
| max_distance | float64 | 最大距離 |
| preferred_tags | []string | 偏好標籤 |
| budget | string | 預算範圍 |
| created_at | time | 建立時間 |
| updated_at | time | 更新時間 |

## 推薦演算法

推薦演算法位於 `internal/services/recommendation.go`，計算方式：

```text
總分 = (天氣適合度 × 天氣權重) + (距離適合度 × 距離權重) + (偏好適合度 × 偏好權重)
```

- **天氣適合度**：根據當前天氣和景點特性計算（0.0 - 1.0）
- **距離適合度**：距離越近分數越高（0.0 - 1.0）
- **偏好適合度**：景點標籤與使用者偏好的匹配度（0.0 - 1.0）

## 測試

### 單元測試

```bash
# 執行所有測試
go test ./...

# 執行 utils 套件測試
go test ./pkg/utils -v
```

### 測試覆蓋率

```bash
# 產生覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### WebSocket 測試

專案提供了一個 WebSocket 測試頁面,位於 `examples/websocket-test.html`。

**使用方式**:

1. 啟動 Tour Server:

   ```bash
   # Windows
   runTour.bat

   # Linux/Mac
   go run -ldflags "-X main.SERVICE_NAME=tour_server -X main.SERVICE_VERSION=0.0.1-dev -X main.SERVICE_ENV=dev" cmd/tour/main.go
   ```

2. 使用瀏覽器開啟測試頁面:

   ```bash
   # 直接在瀏覽器中開啟檔案
   # Windows
   start examples/websocket-test.html

   # macOS
   open examples/websocket-test.html

   # Linux
   xdg-open examples/websocket-test.html
   ```

3. 在測試頁面中:
   - 設定伺服器位址(預設: `ws://localhost:8080/ws`)
   - 輸入客戶端 ID(可選)
   - 點擊「連線」按鈕建立 WebSocket 連線
   - 輸入訊息類型和內容後點擊「發送訊息」
   - 觀察訊息收發狀態

**測試功能**:

- ✅ WebSocket 連線建立和斷線
- ✅ 訊息發送和接收
- ✅ 廣播訊息(向所有客戶端發送)
- ✅ 點對點訊息(向特定客戶端發送)
- ✅ 連線狀態監控
- ✅ 客戶端數量統計
- ✅ 訊息格式驗證

## 擴展伺服器類型

本專案採用 Server 介面設計，支援同時運行多種類型的伺服器。

### 新增伺服器類型

如果需要新增其他類型的伺服器（例如：gRPC、WebSocket），請遵循以下步驟：

1. 在 `internal/server/` 目錄下建立新檔案（例如：`grpc.go`）

1. 實作 Server 介面：

   ```go
   package server

   import (
       "github.com/andy2kuo/TourHelper/internal/logger"
   )

   type GRPCServer struct {
       config *Options
       // ... 其他欄位
   }

   func NewGRPCServer(opts *Options) *GRPCServer {
       return &GRPCServer{
           config: opts,
       }
   }

   func (s *GRPCServer) Start() error {
       logger.Info("gRPC 伺服器啟動")
       // 實作啟動邏輯
       return nil
   }

   func (s *GRPCServer) Stop() error {
       logger.Info("gRPC 伺服器關閉")
       // 實作關閉邏輯
       return nil
   }

   func (s *GRPCServer) Name() string {
       return "gRPC Server"
   }
   ```

1. 在 `cmd/backend/main.go` 或 `cmd/frontend/main.go` 中啟動新的伺服器，或建立新的 main.go：

   ```go
   // 範例：在新的 cmd/grpc/main.go 中
   package main

   import (
       "os"
       "os/signal"
       "syscall"
       "github.com/andy2kuo/TourHelper/service/internal/server"
   )

   func main() {
       // 建立 gRPC 伺服器
       grpcServer := server.NewGRPCServer(opts)
       grpcServer.Start()

       // 等待中斷信號
       quit := make(chan os.Signal, 1)
       signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
       <-quit

       // 關閉伺服器
       grpcServer.Stop()
   }
   ```

   或者在現有的伺服器中同時運行多個伺服器：

   ```go
   // 在 cmd/backend/main.go 中
   // 建立 Backend HTTP 伺服器
   backendServer := server.NewBackendServer(opts)
   backendServer.Start()

   // 建立 gRPC 伺服器
   grpcServer := server.NewGRPCServer(opts)
   grpcServer.Start()

   // 等待中斷信號
   quit := make(chan os.Signal, 1)
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit

   // 關閉所有伺服器
   backendServer.Stop()
   grpcServer.Stop()
   ```

### 優雅關閉

所有伺服器都支援優雅關閉（Graceful Shutdown）：

- Tour HTTP/WebSocket Server：等待現有請求完成（最多 5 秒）
- Lobby HTTP Server：等待現有請求完成（最多 5 秒）
- Backend Admin HTTP Server：等待現有請求完成（最多 5 秒）
- 其他伺服器：實作各自的優雅關閉邏輯

## Redis 整合

### Redis 用途

本專案使用 Redis 作為各伺服器間的狀態交換中心：

1. **會員狀態管理**：Lobby Server 將會員登入狀態寫入 Redis
2. **伺服器狀態同步**：Tour Server 定時將伺服器狀態寫入 Redis
3. **狀態監控**：Lobby Server 和 Backend Server 從 Redis 讀取狀態
4. **Session 管理**：JWT Token 驗證和 Session 資料儲存
5. **快取**：景點資料、推薦結果快取

### Redis Key 設計

```text
member:{memberID}:status      - 會員狀態
tour:server:status            - Tour Server 狀態
tour:destinations             - 景點快取
session:{token}               - Session 資料
cache:recommendation:{params} - 推薦結果快取
```

### Redis 安裝

```bash
# Docker 方式（推薦）
docker run -d -p 6379:6379 --name redis redis:latest

# 或使用持久化
docker run -d -p 6379:6379 --name redis \
  -v redis_data:/data \
  redis:latest redis-server --appendonly yes
```

### Redis 狀態交換範例

#### Lobby Server 更新會員狀態到 Redis

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

#### Lobby Server 定時取得 Tour 狀態

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

#### Tour Server 更新狀態到 Redis

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

### Redis 管理工具

#### 使用 redis-cli 檢查狀態

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

#### 使用 Redis Commander (Web UI)

```bash
# 使用 Docker 啟動
docker run -d -p 8083:8081 --name redis-commander \
  --env REDIS_HOSTS=local:host.docker.internal:6379 \
  rediscommander/redis-commander

# 瀏覽器開啟
# http://localhost:8083
```

### 待實作功能

- [ ] Redis 客戶端連線池
- [ ] 狀態寫入/讀取方法
- [ ] Redis Pub/Sub 即時通知
- [ ] 快取策略實作

## 部署

（待補充：Docker、Docker Compose 等部署方式）

## 日誌管理

### 日誌路徑

應用程式會根據環境自動選擇日誌檔案位置，檔案名稱為 `{SERVICE_NAME}.log`：

**開發環境 (`dev`)**：

- 路徑：`./log/tour_helper.log`
- 相對於專案執行目錄

**生產環境 (Linux/Unix)**：

- 路徑：`/var/log/tour_helper/tour_helper.log`
- 需要適當的檔案系統權限

**生產環境 (Windows)**：

- 路徑：`C:\ProgramData\tour_helper\log\tour_helper.log`
- 需要適當的檔案系統權限

### 日誌 Rotation 設定

日誌檔案會自動 rotation，避免單一檔案過大。可在 `configs/config.yaml` 中設定：

```yaml
log:
  maxsize: 100      # 單一日誌檔案最大大小（MB）
  maxbackups: 3     # 保留的舊日誌檔案數量
  maxage: 28        # 保留的天數
  compress: true    # 是否壓縮舊日誌檔案
```

預設值：

- MaxSize: 100 MB
- MaxBackups: 3 個檔案
- MaxAge: 28 天
- Compress: true

### 日誌級別

- **Debug**：開發環境預設，包含詳細除錯資訊
- **Info**：生產環境預設，一般資訊記錄
- **Warn**：警告訊息
- **Error**：錯誤訊息
- **Fatal**：嚴重錯誤，記錄後程式會終止

### 使用方式

```go
import "github.com/andy2kuo/TourHelper/internal/logger"

// 基本日誌
logger.Info("應用程式啟動")
logger.Infof("伺服器啟動於 %s", addr)

// 帶欄位的日誌
logger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
}).Info("使用者登入")

// 錯誤日誌
logger.Error("發生錯誤")
logger.Errorf("連線失敗: %v", err)
```

## 疑難排解

### 資料庫連線錯誤

確認：

1. 資料庫服務是否執行
2. `configs/config.yaml` 中的 DSN 是否正確
3. 資料庫使用者權限是否足夠

### Bot Webhook 無法接收訊息

確認：

1. Webhook URL 是否設定正確
2. 伺服器是否使用 HTTPS（Line/Telegram 要求）
3. Channel Secret 和 Token 是否正確

### 日誌檔案無法寫入

確認：

1. 日誌目錄是否存在且有寫入權限
2. Linux：檢查 `/var/log/tour_helper/` 權限
3. Windows：檢查 `C:\ProgramData\tour_helper\log\` 權限

## 相關連結

- [Gin Web Framework](https://gin-gonic.com/)
- [GORM 文件](https://gorm.io/)
- [Viper 文件](https://github.com/spf13/viper)
- [Line Messaging API](https://developers.line.biz/en/docs/messaging-api/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
