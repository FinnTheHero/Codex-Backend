package common

import (
	"Codex-Backend/api/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string) (string, error) {

	key, err := GetEnvVariable("JWT_SIGN_KEY")
	if err != nil {
		return "", err
	}
	signKey := []byte(key)

	expirationTime := time.Now().Add(time.Hour * 24)

	calims := &domain.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, calims)

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", &Error{Err: errors.New("Token Service Error - Generate Token String: " + err.Error())}
	}

	return tokenString, nil
}
