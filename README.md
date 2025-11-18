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

### 後端 (Backend)

- **程式語言**：Go 1.25+
- **Web 框架**：Gin Web Framework
- **資料庫**：支援 SQLite / MySQL / PostgreSQL
- **ORM**：GORM
- **配置管理**：Viper
- **Bot 整合**：Line Bot SDK、Telegram Bot API
- **測試框架**：Go testing、Testify

### 前端 (Vue)

- **框架**：Vue.js 3.5+
- **建置工具**：Vite 7.x
- **UI 框架**：Element Plus 2.x
- **狀態管理**：Pinia
- **圖標庫**：Element Plus Icons

## 專案結構

```text
TourHelper/
├── service/                 # 後端 Go 專案
│   ├── cmd/                # 主程式進入點
│   │   ├── backend/        # 後端 API 伺服器
│   │   │   └── main.go
│   │   └── frontend/       # 前端伺服器
│   │       └── main.go
│   ├── internal/           # 私有應用程式碼
│   │   ├── config/         # 設定管理
│   │   ├── server/         # 伺服器實作（backend/frontend）
│   │   ├── services/       # 業務邏輯層
│   │   ├── dao/            # 資料存取層
│   │   ├── models/         # 資料模型層
│   │   ├── database/       # 資料庫管理
│   │   ├── logger/         # 日誌管理
│   │   └── bot/            # Bot 整合（Line/Telegram）
│   ├── pkg/                # 可重用的公開函式庫
│   ├── configs/            # 設定檔
│   └── README.md           # 後端說明文件
├── vue/                     # 前端 Vue.js 專案
│   ├── src/                # 原始碼
│   ├── public/             # 靜態資源
│   └── README.md           # 前端說明文件
├── .claude/                 # Claude Code 設定
├── CLAUDE.md                # Claude 專案說明
└── README.md                # 本檔案
```

詳細的專案結構請參考：

- [後端專案結構](service/README.md#專案結構)
- [前端專案結構](vue/README.md#專案結構)

## 快速開始

### 環境需求

- Go 1.25 或更高版本
- Node.js 18+ (前端開發用)
- （可選）MySQL 或 PostgreSQL

### 後端安裝與執行

1. Clone 專案

   ```bash
   git clone https://github.com/yourusername/TourHelper.git
   cd TourHelper
   ```

1. 進入後端目錄並安裝依賴

   ```bash
   cd service
   go mod download
   ```

1. 設定環境變數

   ```bash
   cp configs/config.example.yaml configs/config.yaml
   # 編輯 config.yaml，填入必要的 API 金鑰
   ```

1. 執行程式

   ```bash
   # 執行後端 API 伺服器
   go run cmd/backend/main.go

   # 執行前端伺服器
   go run cmd/frontend/main.go

   # 或建置後執行
   go build -o backend cmd/backend/main.go
   go build -o frontend cmd/frontend/main.go
   ./backend
   ./frontend
   ```

詳細的後端開發指南請參考 [service/README.md](service/README.md)

### 前端安裝與執行

1. 進入前端目錄並安裝依賴

   ```bash
   cd vue
   npm install
   ```

1. 執行開發伺服器

   ```bash
   npm run dev
   ```

1. 建置生產版本

   ```bash
   npm run build
   ```

詳細的前端開發指南請參考 [vue/README.md](vue/README.md)

## 設定說明

### Line Bot 設定

1. 在 [Line Developers](https://developers.line.biz/) 建立 Messaging API Channel
1. 取得 Channel Secret 和 Channel Access Token
1. 在 `service/configs/config.yaml` 中設定：

   ```yaml
   line:
     enabled: true
     channel_secret: YOUR_CHANNEL_SECRET
     channel_access_token: YOUR_CHANNEL_ACCESS_TOKEN
   ```

### Telegram Bot 設定

1. 透過 [@BotFather](https://t.me/botfather) 建立 Bot
1. 取得 Bot Token
1. 在 `service/configs/config.yaml` 中設定：

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

### 資料庫設定

```yaml
database:
  driver: mysql 
  dsn: tourhelper.db
  # MySQL 範例：
  # dsn: user:password@tcp(localhost:3306)/tourhelper?charset=utf8mb4&parseTime=True
  # PostgreSQL 範例：
  # dsn: host=localhost user=user password=password dbname=tourhelper port=5432 sslmode=disable
```

更多設定說明請參考 [service/README.md#設定說明](service/README.md#設定說明)

## 開發

### 後端開發

請參考 [service/README.md](service/README.md) 了解：

- 開發指令
- 測試方法
- 程式碼品質工具
- 專案架構

### 前端開發

請參考 [vue/README.md](vue/README.md) 了解：

- 開發指令
- 元件說明
- 狀態管理
- 樣式指南

## 部署

（待補充：Docker、Docker Compose 等部署方式）

## 授權

MIT License

## 貢獻

歡迎提交 Issue 或 Pull Request！

## 相關連結

- [Line Bot 開發文件](https://developers.line.biz/en/docs/messaging-api/)
- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM 文件](https://gorm.io/)
- [Vue.js 文件](https://vuejs.org/)
- [Element Plus 文件](https://element-plus.org/)
- [Pinia 文件](https://pinia.vuejs.org/)
