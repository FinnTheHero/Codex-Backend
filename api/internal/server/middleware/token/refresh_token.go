package token

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (mf *IMTokenCache) AutoRefreshTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		defer ctx.Done()

		path := c.Request.URL.Path
		if path == "/user/refresh" || path == "/user/logout" {
			c.Next()
			return
		}

		// Check if access token is close to expiring
		claims, exists := c.Get("claims")
		if !exists {
			c.Next()
			return
		}

		userClaims, ok := claims.(*domain.Claims)
		if !ok {
			c.Next()
			return
		}

		// Refresh if less than 5 minutes remaining on access token
		timeUntilExpiry := time.Until(userClaims.ExpiresAt.Time)
		if timeUntilExpiry > 5*time.Minute {
			c.Next()
			return
		}

		// Access token expires soon, try to refresh using refresh token
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			// No refresh token available, let it expire naturally
			c.Next()
			return
		}

		config := DefaultTokenConfig()

		// Generate new access token using refresh token
		newAccessToken, err := refreshAccessTokenFromString(refreshToken, userClaims.ID, domain.MiddlewareConfig{
			Cache:         mf.cache,
			CacheDuration: 1 * time.Hour,
		}, ctx)
		if err != nil {
			c.Next()
			return
		}

		// Set new access token cookie
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie("access_token", newAccessToken, int(config.AccessTTL.Seconds()), "/", "", true, true)

		// Let frontend know token was refreshed
		c.Header("X-Token-Refreshed", "true")

		c.Next()
	}
}

func refreshAccessTokenFromString(refreshTokenString, expectedUserID string, cacheConfig domain.MiddlewareConfig, ctx context.Context) (string, error) {
	config := DefaultTokenConfig()

	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.SigningKey), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid refresh token claims")
	}

	if claims.Subject != expectedUserID {
		return "", errors.New("refresh token user mismatch")
	}

	// Get user info (from cache or database)
	user, err := service.GetUserByID(claims.ID, ctx)
	if err != nil {
		return "", err
	}

	// Cache user if cache is available
	cacheKey := fmt.Sprintf("user:%s", claims.ID)
	if cacheConfig.Cache != nil {
		cacheConfig.Cache.Set(cacheKey, user, cacheConfig.CacheDuration)
	}

	// Generate new access token
	newAccessToken, _, err := generateAccessToken(user.ID, user.Email, config)
	return newAccessToken, err
}
