package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name:    "default configuration",
			envVars: map[string]string{},
			wantErr: false,
		},
		{
			name: "custom configuration",
			envVars: map[string]string{
				"SERVER_PORT": "9090",
				"GIN_MODE":    "release",
				"DB_HOST":     "testdb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for k, v := range tt.envVars {
				os.Setenv(k, v)
				defer os.Unsetenv(k)
			}

			cfg, err := Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Validate configuration
				if cfg == nil {
					t.Error("Load() returned nil config")
					return
				}

				if cfg.Server.Port == "" {
					t.Error("Server port is empty")
				}

				// Check custom values
				if port, ok := tt.envVars["SERVER_PORT"]; ok {
					if cfg.Server.Port != port {
						t.Errorf("Server.Port = %v, want %v", cfg.Server.Port, port)
					}
				}

				if mode, ok := tt.envVars["GIN_MODE"]; ok {
					if cfg.Server.Mode != mode {
						t.Errorf("Server.Mode = %v, want %v", cfg.Server.Mode, mode)
					}
				}
			}
		})
	}
}

func TestDatabaseConfig_GetDSN(t *testing.T) {
	cfg := &DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		DBName:   "testdb",
		SSLMode:  "disable",
	}

	dsn := cfg.GetDSN()
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"

	if dsn != expected {
		t.Errorf("GetDSN() = %v, want %v", dsn, expected)
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_VAR",
			defaultValue: "default",
			envValue:     "custom",
			want:         "custom",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_VAR",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := getEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
