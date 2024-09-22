package auth_services

import (
	"Codex-Backend/api/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {

	signKey := []byte(os.Getenv("JWT_SIGN_KEY"))

	expirationTime := time.Now().Add(time.Hour * 24)

	calims := &models.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, calims)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
