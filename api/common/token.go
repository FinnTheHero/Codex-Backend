package common

import (
	"Codex-Backend/api/internal/domain"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokenPair(ID, email string, config domain.TokenConfig) (*domain.TokenPair, error) {
	if ID == "" {
		return nil, &Error{Err: errors.New("user ID cannot be empty")}
	}
	if email == "" {
		return nil, &Error{Err: errors.New("email cannot be empty")}
	}
	if config.SigningKey == "" {
		return nil, &Error{Err: errors.New("signing key not configured")}
	}

	// Generate access token
	accessToken, expiresAt, err := generateAccessToken(ID, email, config)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := generateRefreshToken(ID, config)
	if err != nil {
		return nil, err
	}

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
	}, nil
}

// GenerateAccessToken creates a new access token (for refresh scenarios)
func GenerateAccessToken(ID, email string, config domain.TokenConfig) (string, time.Time, error) {
	return generateAccessToken(ID, email, config)
}

// generateAccessToken creates the actual access token
func generateAccessToken(ID, email string, config domain.TokenConfig) (string, time.Time, error) {
	now := time.Now()
	expirationTime := now.Add(config.AccessTTL)

	// Generate unique token ID for revocation capability
	jti, err := generateJTI()
	if err != nil {
		return "", time.Time{}, &Error{Err: errors.New("failed to generate token ID: " + err.Error())}
	}

	claims := &domain.Claims{
		ID:    ID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Subject:   ID,
			Audience:  jwt.ClaimStrings{config.Audience},
			Issuer:    config.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SigningKey))
	if err != nil {
		return "", time.Time{}, &Error{Err: errors.New("token signing failed: " + err.Error())}
	}

	return tokenString, expirationTime, nil
}

func generateRefreshToken(userID string, config domain.TokenConfig) (string, error) {
	now := time.Now()
	expirationTime := now.Add(config.RefreshTTL)

	jti, err := generateJTI()
	if err != nil {
		return "", &Error{Err: errors.New("failed to generate refresh token ID: " + err.Error())}
	}

	// Refresh tokens have minimal claims
	claims := &jwt.RegisteredClaims{
		ID:        jti,
		Subject:   userID,
		Audience:  jwt.ClaimStrings{config.Audience},
		Issuer:    config.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		NotBefore: jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.SigningKey))
	if err != nil {
		return "", &Error{Err: errors.New("refresh token signing failed: " + err.Error())}
	}

	return tokenString, nil
}

func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func DefaultTokenConfig() domain.TokenConfig {
	return domain.TokenConfig{
		SigningKey: GetEnvVariable("JWT_SIGN_KEY"),
		AccessTTL:  30 * time.Minute,
		RefreshTTL: 7 * 24 * time.Hour,
		Issuer:     GetEnvVariable("JWT_ISSUER"),
		Audience:   GetEnvVariable("JWT_AUDIENCE"),
	}
}
