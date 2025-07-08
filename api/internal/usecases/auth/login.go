package auth_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"errors"
)

func LoginUser(credentials domain.Credentials) (string, domain.User, error) {
	var user domain.User

	// Check if user exists
	user, err := repository.GetUser(credentials.Email)
	if err != nil {
		return "", user, err
	}

	if user.Email == "" {
		return "", user, errors.New("User not found")
	}

	// Check if password is correct
	err = VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return "", user, errors.New("Incorrect password")
	}

	// Generate token
	token, err := GenerateToken(user.Email)
	if err != nil {
		return "", user, errors.New("Error generating token: " + err.Error())
	}

	return token, user, nil
}
