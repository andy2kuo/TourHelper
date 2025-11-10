package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test with default values
	cfg := Load()

	if cfg == nil {
		t.Fatal("Load() returned nil")
	}

	if cfg.ServerPort == "" {
		t.Error("ServerPort should not be empty")
	}

	if cfg.JWTSecret == "" {
		t.Error("JWTSecret should not be empty")
	}

	if cfg.AllowedOrigins == nil {
		t.Error("AllowedOrigins should not be nil")
	}
}

func TestLoadWithEnvVars(t *testing.T) {
	// Set environment variables
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("GOOGLE_MAPS_API_KEY", "test-api-key")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("GOOGLE_MAPS_API_KEY")
	}()

	cfg := Load()

	if cfg.ServerPort != "9090" {
		t.Errorf("Expected ServerPort to be 9090, got %s", cfg.ServerPort)
	}

	if cfg.JWTSecret != "test-secret" {
		t.Errorf("Expected JWTSecret to be test-secret, got %s", cfg.JWTSecret)
	}

	if cfg.GoogleMapsAPIKey != "test-api-key" {
		t.Errorf("Expected GoogleMapsAPIKey to be test-api-key, got %s", cfg.GoogleMapsAPIKey)
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test-value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default")
	if value != "test-value" {
		t.Errorf("Expected test-value, got %s", value)
	}

	// Test with non-existing env var
	value = getEnv("NON_EXISTING_VAR", "default")
	if value != "default" {
		t.Errorf("Expected default, got %s", value)
	}
}
