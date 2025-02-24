package auth_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"errors"
)

func (s *AuthService) LoginUser(credentials domain.Credentials) (string, domain.User, error) {
	var user domain.User

	// Check if user exists
	result, err := repository.GetUser(credentials.Email)
	if err != nil {
		return "", user, err
	}

	userDTO, ok := result.(domain.UserDTO)
	if !ok {
		return "", user, errors.New("Type assertion failed")
	}

	user = userDTO.User

	if userDTO.Email == "" {
		return "", user, errors.New("User not found")
	}

	// Check if password is correct
	err = VerifyPassword(userDTO.User.Password, credentials.Password)
	if err != nil {
		return "", user, errors.New("Incorrect password")
	}

	// Generate token
	token, err := GenerateToken(userDTO.User.Email)
	if err != nil {
		return "", user, errors.New("Error generating token: " + err.Error())
	}

	return token, user, nil
}
