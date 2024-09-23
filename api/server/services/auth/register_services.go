package auth_services

import (
	aws_services "Codex-Backend/api/aws/services"
	"Codex-Backend/api/models"
	"errors"
	"time"
)

func (s *AuthService) RegisterUser(user models.User) error {

	// Check if user already exists
	result, err := aws_services.GetUser(user.Email)
	if err != nil {
		if err.Error() != "User not found" {
			return err
		}
	}

	userDTO, _ := result.(models.UserDTO)

	if userDTO.Email == user.Email || userDTO.Email != "" {
		return errors.New("Email already in use")
	}

	// Hash password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return errors.New("Error hashing password: " + err.Error())
	}
	user.Password = string(hashedPassword)

	// Set user defaults
	user.Type = "user"
	user.Created_at = time.Now().Format("2006-01-02 15:04:05")
	user.Updated_at = user.Created_at

	// Create new user
	err = aws_services.CreateUser(user)
	if err != nil {
		return errors.New("Error creating user: " + err.Error())
	}

	return nil
}
