# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

請保持繁體中文回覆。

## README.md 更新指引

每次更新時，請同步更新以下 README.md 中的相關內容，以確保說明文件與專案狀態一致：
  
### ./README.md

專案說明文件，包含專案概述、快速開始指南、開發指令等資訊。

### ./backend/README.md

後端 Go 專案說明文件，包含後端架構、API 文件、資料庫設計等資訊。

### ./vue/README.md

前端 Vue.js 專案說明文件，包含前端架構、元件說明、狀態管理等資訊。

## Project Overview

TourHelper is a Go application designed to automatically suggest travel destinations. The README states: "自動幫我想可以去哪旅遊" (Automatically help me think of where to travel).

這個專案的用途主要是用來自動推薦我現在可以去哪裡玩，需要根據當時位置、天氣、距離等條件來做推薦。專案中包含了與 Line 和 Telegram 的機器人整合，讓使用者可以透過這些平台獲取旅遊建議。

主要呈現介面有：

* Line 機器人介面
* Telegram 機器人介面
* 網頁介面

考慮到有可能會有即時互動需求，專案中也可能包含 WebSocket 功能來提供即時更新。

可能會有不同用戶的需求，因此專案中也可能包含用戶管理和偏好設定功能。

## Language and Framework

* Go (Golang) 1.25+
* Web frameworks/libraries: Gin
* Vue.js for frontend
* Line Bot SDK for Line integration
* Telegram Bot API for Telegram integration
* Database: MySQL
* WebSocket for real-time features
* Other libraries: GORM for ORM, Viper for configuration management
* Testing: Go's built-in testing package, Testify
* 運行環境: Docker (if applicable)
* Logger: Logrus

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

請參考 ./backend/README.md 和 ./vue/README.md 了解後端和前端的詳細專案結構。

### 後端分層架構

本專案後端採用經典的分層架構設計：

```text
外部請求（HTTP / Bot / gRPC 等）
    ↓
[Handlers] ← server/handlers.go（處理請求/回應、參數驗證）
    ↓
[Services] ← services/（業務邏輯處理）
    ↓
[DAO] ← dao/（資料庫 CRUD 操作）
    ↓
[Models] ← models/（資料結構定義）
    ↓
資料庫
```

### 重要套件說明

* **cmd/tourhelper/main.go**: 應用程式進入點，初始化設定、建立 server、優雅關閉
* **internal/config/**: 設定管理，使用 Viper 載入 YAML 設定檔和環境變數
* **internal/server/**: 伺服器實作，包含 Server 介面、HTTP 伺服器和 Handlers
* **internal/services/**: 業務邏輯層，處理核心業務邏輯（使用單例模式）
* **internal/dao/**: 資料存取層，封裝所有資料庫 CRUD 操作（使用單例模式）
* **internal/models/**: 資料模型層，定義資料結構和 GORM 標籤
* **internal/logger/**: 日誌管理，使用 Logrus + Lumberjack
* **internal/bot/**: Bot 整合，包含 Line 和 Telegram Bot 實作

## Environment Configuration

Environment variables should be stored in `.env` file (already in .gitignore). Common variables may include:

* API keys for travel services
* Database connection strings
* External service endpoints
