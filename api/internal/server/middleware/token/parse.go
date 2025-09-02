package token

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Parses and validates the JWT token
func ParseAndValidateJWT(tokenString string) (*domain.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key := cmn.GetEnvVariable("JWT_SIGN_KEY")
		return []byte(key), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*domain.Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate required fields
	if claims.ID == "" {
		return nil, fmt.Errorf("user ID not found in token")
	}
	if claims.Email == "" {
		return nil, fmt.Errorf("email not found in token")
	}

	return claims, nil
}
