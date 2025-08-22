package domain

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Type  string `json:"type"`
	jwt.RegisteredClaims
}

type TokenConfig struct {
	SigningKey string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string // e.g., "your-app-name"
	Audience   string // e.g., "your-app-users"
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// TokenCache interface for optional token caching
type TokenCache interface {
	Get(key string) (any, bool)
	Set(key string, value any, duration time.Duration)
	Delete(key string)
}

type LookupUser struct {
	Cache         TokenCache
	CacheDuration time.Duration
}

type InMemoryCache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value     any
	expiresAt time.Time
}
