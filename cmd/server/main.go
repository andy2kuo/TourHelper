package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/andy2kuo/TourHelper/config"
	"github.com/andy2kuo/TourHelper/internal/auth"
	"github.com/andy2kuo/TourHelper/internal/handlers"
	"github.com/andy2kuo/TourHelper/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize stores and hubs
	authStore := auth.NewStore()
	wsHub := websocket.NewHub()

	// Start WebSocket hub
	go wsHub.Run()

	// Create handler
	handler := handlers.NewHandler(cfg, authStore, wsHub)

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/login", handler.HandleLogin)
	mux.HandleFunc("/api/config", handler.HandleGetConfig)
	mux.HandleFunc("/api/health", handler.HandleHealth)
	mux.HandleFunc("/ws", handler.HandleWebSocket)

	// Serve static files
	webDir := filepath.Join(".", "web")
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Enable CORS
	httpHandler := handlers.EnableCORS(mux)

	// Start server
	addr := ":" + cfg.ServerPort
	log.Printf("Server starting on %s", addr)
	log.Printf("Web client available at http://localhost%s", addr)
	log.Printf("Default credentials:")
	log.Printf("  Web: webuser / password123")
	log.Printf("  Line: lineuser / password123")
	log.Printf("  Telegram: telegramuser / password123")

	if err := http.ListenAndServe(addr, httpHandler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
