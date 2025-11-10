# TourHelper

自動幫我想可以去哪旅遊

TourHelper 是一個智慧旅遊推薦系統，根據使用者的當前位置、即時天氣、個人偏好等因素，自動推薦適合的旅遊景點。

## 功能特色

- 🌍 **智慧推薦**：根據位置、天氣、距離等多重因素智慧推薦景點
- 💬 **多平台支援**：支援 Line Bot、Telegram Bot 和網頁介面
- 🌦️ **天氣整合**：即時天氣資訊，推薦適合當前天氣的活動
- 📍 **距離計算**：精確計算距離和預估旅行時間
- ⭐ **個人化偏好**：可設定距離範圍、景點類型、預算等偏好
- 📊 **歷史記錄**：記錄搜尋和造訪歷史

## 技術架構

- **後端**：Go 1.25+ with Gin Web Framework
- **前端**：Vue.js（開發中）
- **資料庫**：支援 SQLite / MySQL / PostgreSQL
- **Bot 整合**：Line Bot SDK、Telegram Bot API
- **配置管理**：Viper
- **ORM**：GORM

## 快速開始

### 環境需求

- Go 1.25 或更高版本
- （可選）MySQL 或 PostgreSQL

### 安裝步驟

1. Clone 專案

```bash
git clone https://github.com/yourusername/TourHelper.git
cd TourHelper
```

2. 安裝依賴

```bash
go mod download
```

3. 設定環境變數

```bash
cp .env.example .env
cp configs/config.example.yaml configs/config.yaml
# 編輯 .env 或 config.yaml，填入必要的 API 金鑰
```

4. 執行程式

```bash
# 開發模式
go run cmd/tourhelper/main.go

# 或建置後執行
go build -o tourhelper cmd/tourhelper/main.go
./tourhelper
```

## 設定說明

### Line Bot 設定

1. 在 [Line Developers](https://developers.line.biz/) 建立 Messaging API Channel
2. 取得 Channel Secret 和 Channel Access Token
3. 在 `config.yaml` 或環境變數中設定：

   ```yaml
   line:
     enabled: true
     channel_secret: YOUR_CHANNEL_SECRET
     channel_access_token: YOUR_CHANNEL_ACCESS_TOKEN
   ```

### Telegram Bot 設定

1. 透過 [@BotFather](https://t.me/botfather) 建立 Bot
2. 取得 Bot Token
3. 在 `config.yaml` 或環境變數中設定：

   ```yaml
   telegram:
     enabled: true
     token: YOUR_BOT_TOKEN
   ```

### 天氣 API 設定

支援多種天氣服務供應商（OpenWeatherMap、WeatherAPI 等）：

```yaml
weather:
  api_key: YOUR_API_KEY
  provider: openweathermap
```

### 地圖 API 設定

支援 Google Maps、HERE、Mapbox 等：

```yaml
maps:
  api_key: YOUR_API_KEY
  provider: google
```

## 開發指令

```bash
# 執行測試
go test ./...

# 執行特定測試
go test -run TestFunctionName ./path/to/package

# 格式化程式碼
go fmt ./...

# 執行 linter
golangci-lint run

# 建置
go build -o tourhelper cmd/tourhelper/main.go
```

## 專案結構

```text
TourHelper/
├── cmd/
│   └── tourhelper/      # 主程式進入點
├── internal/
│   ├── config/          # 設定管理
│   ├── models/          # 資料模型
│   ├── handlers/        # HTTP 處理器
│   ├── services/        # 業務邏輯服務
│   ├── bot/
│   │   ├── line/        # Line Bot 整合
│   │   └── telegram/    # Telegram Bot 整合
│   └── middleware/      # 中介軟體
├── pkg/
│   └── utils/           # 工具函式
├── configs/             # 設定檔
├── web/                 # 前端檔案
└── test/                # 測試檔案
```

## API 端點

### REST API

- `GET /health` - 健康檢查
- `GET /api/v1/recommendations` - 取得旅遊推薦
- `POST /api/v1/user/preferences` - 更新使用者偏好
- `GET /api/v1/user/preferences` - 取得使用者偏好

### Webhook

- `POST /webhook/line` - Line Bot webhook
- `POST /webhook/telegram` - Telegram Bot webhook

### WebSocket

- `GET /ws` - WebSocket 連線（開發中）

## 授權

MIT License

## 貢獻

歡迎提交 Issue 或 Pull Request！
