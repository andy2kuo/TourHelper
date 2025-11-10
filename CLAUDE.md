# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

請保持繁體中文回覆。

每次進行調整後，請自動更新README.md的說明檔案，確保README.md包含以下內容：

- 專案簡介
- 使用的程式語言和框架
- 主要的開發指令（如建置、執行、測試等）
- 專案結構說明
- 環境設定說明
- 其他相關資訊

## Project Overview

TourHelper is a Go application designed to automatically suggest travel destinations. The README states: "自動幫我想可以去哪旅遊" (Automatically help me think of where to travel).

這個專案的用途主要是用來自動推薦我現在可以去哪裡玩，需要根據當時位置、天氣、距離等條件來做推薦。專案中包含了與 Line 和 Telegram 的機器人整合，讓使用者可以透過這些平台獲取旅遊建議。

主要呈現介面有：

- Line 機器人介面
- Telegram 機器人介面
- 網頁介面

考慮到有可能會有即時互動需求，專案中也可能包含 WebSocket 功能來提供即時更新。

可能會有不同用戶的需求，因此專案中也可能包含用戶管理和偏好設定功能。

## Language and Framework

- Go (Golang) 1.25+
- Web frameworks/libraries: Gin
- Vue.js for frontend (if applicable)
- Line Bot SDK for Line integration
- Telegram Bot API for Telegram integration
- Database: MySQL or PostgreSQL (if applicable)
- WebSocket for real-time features (if applicable)
- Other libraries: GORM for ORM, Viper for configuration management
- Testing: Go's built-in testing package, Testify for assertions
- 運行環境: Docker (if applicable)
- Logger: Logrus or Zap (if applicable)

## Development Commands

### Project Initialization

```bash
# Initialize Go module (if not already done)
go mod init github.com/yourusername/TourHelper

# Download dependencies
go mod download

# Tidy up dependencies
go mod tidy
```

### Building

```bash
# Build the application
go build -o tourhelper

# Build with specific output location
go build -o bin/tourhelper ./cmd/tourhelper
```

### Running

```bash
# Run directly without building
go run main.go

# Run with arguments
go run main.go [args]
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run a specific test
go test -run TestFunctionName ./path/to/package

# Run tests in a specific package
go test ./internal/packagename
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter (requires golangci-lint)
golangci-lint run

# Run go vet
go vet ./...

# Check for common mistakes
go vet ./...
```

## Project Structure

目前專案已建立以下結構：

```text
TourHelper/
├── cmd/tourhelper/          # 主程式進入點
│   └── main.go             # 應用程式啟動檔案，初始化 Gin server、載入設定、註冊路由
├── internal/               # 私有應用程式碼
│   ├── config/            # 設定管理
│   │   └── config.go      # 使用 Viper 管理設定，支援 YAML 和環境變數
│   ├── models/            # 資料模型
│   │   ├── models.go      # 定義 User, Destination, Tag 等資料結構
│   │   └── database.go    # 資料庫初始化、遷移和範例資料填充
│   ├── handlers/          # HTTP 請求處理器
│   │   └── handlers.go    # 處理 API 請求（推薦、偏好設定等）
│   ├── services/          # 業務邏輯服務
│   │   ├── recommendation.go  # 推薦演算法核心邏輯
│   │   └── weather.go         # 天氣資訊服務
│   └── bot/              # Bot 整合
│       ├── line/         # Line Bot 實作
│       │   └── line.go   # 處理 Line webhook、訊息回覆
│       └── telegram/     # Telegram Bot 實作
│           └── telegram.go  # 處理 Telegram webhook、指令處理
├── pkg/                  # 可重用的公開函式庫
│   └── utils/           # 工具函式
│       ├── utils.go     # 距離計算、時間估算等工具
│       └── utils_test.go # 工具函式測試
├── configs/             # 設定檔
│   └── config.example.yaml  # 設定檔範例
├── .env.example         # 環境變數範例
├── .air.toml           # Air 自動重載設定
├── Makefile            # 開發指令快捷方式
└── go.mod              # Go 模組依賴
```

### 重要檔案說明

- **main.go**: 應用程式進入點，初始化設定、建立 Gin router、註冊所有路由和中介軟體
- **config.go**: 使用 Viper 載入 YAML 設定檔和環境變數
- **models.go**: 定義所有資料模型，包含 GORM 標籤
- **database.go**: 資料庫連線初始化、自動遷移和範例資料填充功能
- **handlers.go**: HTTP API 端點處理器
- **recommendation.go**: 推薦演算法，根據位置、天氣、距離計算適合度評分
- **line.go / telegram.go**: Bot 整合實作，處理訊息和位置資訊

## Environment Configuration

Environment variables should be stored in `.env` file (already in .gitignore). Common variables may include:

- API keys for travel services
- Database connection strings
- External service endpoints
