package token

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/service"
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func refreshAccessTokenFromString(refreshTokenString, expectedUserID string, cacheConfig domain.LookupUser, ctx context.Context) (string, error) {
	config := DefaultTokenConfig()

	// Parse refresh token
	token, err := jwt.ParseWithClaims(refreshTokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
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
	newAccessToken, _, err := generateAccessToken(user, config)
	return newAccessToken, err
}
