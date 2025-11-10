package auth

import (
	"testing"
	"time"

	"github.com/andy2kuo/TourHelper/internal/models"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Fatal("NewStore returned nil")
	}

	// Check that default users are created
	user, err := store.GetUser("webuser")
	if err != nil {
		t.Errorf("Expected webuser to exist, got error: %v", err)
	}
	if user == nil {
		t.Error("Expected user to be non-nil")
	}
}

func TestAuthenticate(t *testing.T) {
	store := NewStore()

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{
			name:     "valid web user",
			username: "webuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "valid line user",
			username: "lineuser",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "invalid password",
			username: "webuser",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "non-existent user",
			username: "nonexistent",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := store.Authenticate(tt.username, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && user == nil {
				t.Error("Expected user to be non-nil for valid credentials")
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	user := &models.User{
		ID:        "1",
		Username:  "testuser",
		Platform:  "web",
		CreatedAt: time.Now(),
	}

	secret := "test-secret"
	token, err := GenerateToken(user, secret)

	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
		return
	}

	if token == "" {
		t.Error("Expected token to be non-empty")
	}
}

func TestValidateToken(t *testing.T) {
	user := &models.User{
		ID:        "1",
		Username:  "testuser",
		Platform:  "web",
		CreatedAt: time.Now(),
	}

	secret := "test-secret"
	token, err := GenerateToken(user, secret)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Test valid token
	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
		return
	}

	if claims.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, claims.Username)
	}

	if claims.Platform != user.Platform {
		t.Errorf("Expected platform %s, got %s", user.Platform, claims.Platform)
	}

	// Test invalid token
	_, err = ValidateToken("invalid-token", secret)
	if err == nil {
		t.Error("Expected error for invalid token")
	}

	// Test wrong secret
	_, err = ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Error("Expected error for wrong secret")
	}
}

func TestHashPassword(t *testing.T) {
	password := "testpassword"
	hash1 := hashPassword(password)
	hash2 := hashPassword(password)

	if hash1 != hash2 {
		t.Error("Same password should produce same hash")
	}

	differentHash := hashPassword("differentpassword")
	if hash1 == differentHash {
		t.Error("Different passwords should produce different hashes")
	}
}
