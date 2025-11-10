package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort      string
	JWTSecret       string
	GoogleMapsAPIKey string
	AllowedOrigins  []string
}

func Load() *Config {
	return &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		GoogleMapsAPIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		AllowedOrigins:  []string{"*"},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
