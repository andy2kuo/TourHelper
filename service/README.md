# TourHelper Backend

TourHelper 後端服務，使用 Go 語言和 Gin 框架建構，提供旅遊推薦 API、Bot 整合等功能。

## 技術架構

- **程式語言**：Go 1.25+
- **Web 框架**：Gin Web Framework
- **資料庫**：支援 SQLite / MySQL / PostgreSQL
- **ORM**：GORM
- **配置管理**：Viper
- **Bot 整合**：Line Bot SDK、Telegram Bot API
- **測試框架**：Go testing、Testify
- **日誌**：標準 log 套件（可擴展為 Logrus 或 Zap）

## 專案結構

```text
service/
├── cmd/                    # 主程式進入點
│   ├── backend/            # 後端 API 伺服器
│   │   └── main.go         # 後端伺服器啟動檔案，初始化設定、建立 backend server、優雅關閉
│   └── frontend/           # 前端伺服器
│       └── main.go         # 前端伺服器啟動檔案，初始化設定、建立 frontend server、優雅關閉
├── internal/               # 私有應用程式碼
│   ├── config/             # 設定管理
│   │   └── config.go       # 使用 Viper 管理設定，支援 YAML 和環境變數
│   ├── server/             # 伺服器實作
│   │   ├── server.go       # Server 介面定義
│   │   ├── backend/        # 後端 API 伺服器實作
│   │   │   ├── backend-server.go   # Backend 伺服器實作（Gin）
│   │   │   └── backend-handler.go  # Backend API 請求處理器
│   │   └── frontend/       # 前端伺服器實作
│   │       ├── frontend-server.go  # Frontend 伺服器實作（Gin）
│   │       └── frontend-handler.go # Frontend 請求處理器和 Bot webhook
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

## 分層架構

本專案採用經典的分層架構設計，確保關注點分離和程式碼可維護性：

```text
外部請求（HTTP / Bot / gRPC 等）
    ↓
[Handlers] ← server/backend/backend-handler.go 或 server/frontend/frontend-handler.go
    ↓ 處理請求/回應、參數驗證
[Services] ← services/
    ↓ 業務邏輯處理
[DAO] ← dao/
    ↓ 資料庫 CRUD 操作
[Models] ← models/
    ↓ 資料結構定義
資料庫
```

### 各層職責

1. **Handlers（處理器層）**
   - 位置：`internal/server/backend/backend-handler.go` 和 `internal/server/frontend/frontend-handler.go`
   - 職責：處理各種請求和回應、參數解析和驗證、呼叫 Service 層
   - 範例：Backend 的 API handlers、Frontend 的 HealthCheckHandler 和 Bot webhook handlers

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

### cmd/backend/main.go

後端 API 伺服器進入點，負責：

- 載入設定（使用 config.Load）
- 初始化 Logger
- 建立 Backend Server 實例
- 啟動後端 API 伺服器
- 處理優雅關閉（Graceful Shutdown）

### cmd/frontend/main.go

前端伺服器進入點，負責：

- 載入設定（使用 config.Load）
- 初始化 Logger
- 建立 Frontend Server 實例
- 啟動前端伺服器（包含靜態檔案服務和 Bot webhook）
- 處理優雅關閉（Graceful Shutdown）

### internal/server/

伺服器實作模組，支援多種伺服器類型：

- **server.go**：定義 Server 介面，規範所有伺服器類型的行為
  - `Start()`：啟動伺服器
  - `Stop()`：停止伺服器
  - `Name()`：返回伺服器名稱

- **backend/**：後端 API 伺服器實作
  - **backend-server.go**：使用 Gin 框架實作 Backend HTTP 伺服器
    - 註冊後端 API 路由
    - 支援優雅關閉
  - **backend-handler.go**：後端 API 請求處理器
    - 未來可新增各種 API 端點處理器

- **frontend/**：前端伺服器實作
  - **frontend-server.go**：使用 Gin 框架實作 Frontend HTTP 伺服器
    - 註冊前端路由
    - 支援靜態檔案服務
    - 整合 Line 和 Telegram Bot webhook
    - 支援優雅關閉
  - **frontend-handler.go**：前端請求處理器
    - HealthCheckHandler：處理健康檢查請求
    - Bot webhook handlers

**設計理念**：透過 Server 介面和前後端分離架構，可以：

1. 獨立部署和擴展前後端伺服器
2. 輕鬆新增其他類型的伺服器（例如：gRPC Server、WebSocket Server 等）
3. 清晰的職責劃分：Backend 專注於 API 服務，Frontend 專注於頁面和 Bot 整合

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

```bash
# 開發模式（直接執行）
# 執行後端 API 伺服器
go run cmd/backend/main.go

# 執行前端伺服器
go run cmd/frontend/main.go

# 或指定參數
go run cmd/backend/main.go --config configs/config.yaml
go run cmd/frontend/main.go --config configs/config.yaml
```

### 建置

```bash
# 基本建置
go build -o backend cmd/backend/main.go
go build -o frontend cmd/frontend/main.go

# 建置到特定目錄
go build -o bin/backend cmd/backend/main.go
go build -o bin/frontend cmd/frontend/main.go

# 建置並注入版本資訊（使用 -ldflags）
# Backend 伺服器
go build -ldflags "\
  -X main.SERVICE_NAME=tour_helper_backend \
  -X main.SERVICE_ENV=production \
  -X main.SERVICE_VERSION=1.0.0" \
  -o backend cmd/backend/main.go

# Frontend 伺服器
go build -ldflags "\
  -X main.SERVICE_NAME=tour_helper_frontend \
  -X main.SERVICE_ENV=production \
  -X main.SERVICE_VERSION=1.0.0" \
  -o frontend cmd/frontend/main.go

# 開發環境建置
go build -ldflags "\
  -X main.SERVICE_NAME=tour_helper_backend \
  -X main.SERVICE_ENV=dev \
  -X main.SERVICE_VERSION=0.0.1-dev" \
  -o backend cmd/backend/main.go

go build -ldflags "\
  -X main.SERVICE_NAME=tour_helper_frontend \
  -X main.SERVICE_ENV=dev \
  -X main.SERVICE_VERSION=0.0.1-dev" \
  -o frontend cmd/frontend/main.go

# 執行建置的檔案
./backend    # 啟動後端 API 伺服器
./frontend   # 啟動前端伺服器
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
  "status": "ok"
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

- Backend HTTP Server：等待現有請求完成（最多 5 秒）
- Frontend HTTP Server：等待現有請求完成（最多 5 秒）
- 其他伺服器：實作各自的優雅關閉邏輯

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
