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

請參考./backend/README.md 和 ./vue/README.md 了解後端和前端的詳細專案結構。

### 重要檔案說明

* **main.go**: 應用程式進入點，初始化設定、建立 Gin router、註冊所有路由和中介軟體
* **config.go**: 使用 Viper 載入 YAML 設定檔和環境變數
* **models.go**: 定義所有資料模型，包含 GORM 標籤
* **database.go**: 資料庫連線初始化、自動遷移和範例資料填充功能
* **handlers.go**: HTTP API 端點處理器
* **recommendation.go**: 推薦演算法，根據位置、天氣、距離計算適合度評分
* **line.go / telegram.go**: Bot 整合實作，處理訊息和位置資訊

## Environment Configuration

Environment variables should be stored in `.env` file (already in .gitignore). Common variables may include:

* API keys for travel services
* Database connection strings
* External service endpoints
