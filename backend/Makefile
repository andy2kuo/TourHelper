.PHONY: help build run test clean fmt lint install

help: ## 顯示說明
	@echo "可用的指令："
	@echo "  make build   - 建置應用程式"
	@echo "  make run     - 執行應用程式"
	@echo "  make test    - 執行測試"
	@echo "  make fmt     - 格式化程式碼"
	@echo "  make lint    - 執行 linter"
	@echo "  make clean   - 清理建置檔案"
	@echo "  make install - 安裝依賴"

build: ## 建置應用程式
	@echo "建置 TourHelper..."
	go build -o bin/tourhelper cmd/tourhelper/main.go

run: ## 執行應用程式
	@echo "執行 TourHelper..."
	go run cmd/tourhelper/main.go

test: ## 執行所有測試
	@echo "執行測試..."
	go test -v ./...

test-coverage: ## 執行測試並產生覆蓋率報告
	@echo "執行測試並產生覆蓋率報告..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

fmt: ## 格式化程式碼
	@echo "格式化程式碼..."
	go fmt ./...

lint: ## 執行 linter
	@echo "執行 linter..."
	golangci-lint run

clean: ## 清理建置檔案
	@echo "清理建置檔案..."
	rm -rf bin/
	rm -f coverage.out

install: ## 安裝依賴
	@echo "安裝依賴..."
	go mod download
	go mod tidy

dev: ## 開發模式（使用 air 自動重載）
	@echo "啟動開發模式..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "請先安裝 air: go install github.com/cosmtrek/air@latest"; \
	fi
