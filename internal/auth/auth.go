package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/andy2kuo/TourHelper/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
)

// Store is a simple in-memory user store
type Store struct {
	users map[string]*models.User
	mu    sync.RWMutex
}

// NewStore creates a new user store
func NewStore() *Store {
	store := &Store{
		users: make(map[string]*models.User),
	}
	// Add default test users for different platforms
	store.addDefaultUsers()
	return store
}

func (s *Store) addDefaultUsers() {
	// Web user
	s.users["webuser"] = &models.User{
		ID:        "1",
		Username:  "webuser",
		Password:  hashPassword("password123"),
		Email:     "web@example.com",
		Platform:  "web",
		CreatedAt: time.Now(),
	}
	
	// Line user
	s.users["lineuser"] = &models.User{
		ID:        "2",
		Username:  "lineuser",
		Password:  hashPassword("password123"),
		Email:     "line@example.com",
		Platform:  "line",
		CreatedAt: time.Now(),
	}
	
	// Telegram user
	s.users["telegramuser"] = &models.User{
		ID:        "3",
		Username:  "telegramuser",
		Password:  hashPassword("password123"),
		Email:     "telegram@example.com",
		Platform:  "telegram",
		CreatedAt: time.Now(),
	}
}

// Authenticate verifies user credentials
func (s *Store) Authenticate(username, password string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	if user.Password != hashPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// GetUser retrieves a user by username
func (s *Store) GetUser(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Claims represents JWT claims
type Claims struct {
	Username string `json:"username"`
	Platform string `json:"platform"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func GenerateToken(user *models.User, secret string) (string, error) {
	claims := &Claims{
		Username: user.Username,
		Platform: user.Platform,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
