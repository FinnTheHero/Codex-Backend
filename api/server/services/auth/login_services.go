package auth_services

import (
	aws_services "Codex-Backend/api/aws/services"
	"Codex-Backend/api/models"
	"errors"
)

func (s *AuthService) LoginUser(credentials models.Credentials) (string, models.User, error) {
	var user models.User

	// Check if user exists
	result, err := aws_services.GetUser(credentials.Email)
	if err != nil {
		return "", user, err
	}

	userDTO, ok := result.(models.UserDTO)
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
