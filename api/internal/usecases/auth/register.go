package auth_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"errors"
	"time"
)

func (s *AuthService) RegisterUser(credentials domain.NewUser) error {

	var user domain.User
	user.Email = credentials.Email
	user.Username = credentials.Username

	// Check if user already exists
	result, err := repository.GetUser(credentials.Email)
	if err != nil {
		if err.Error() != "User not found" {
			return err
		}
	}

	userDTO, _ := result.(domain.UserDTO)

	if userDTO.Email == credentials.Email || userDTO.Email != "" {
		return errors.New("Email already in use")
	}

	// Hash password
	hashedPassword, err := HashPassword(credentials.Password)
	if err != nil {
		return errors.New("Error hashing password: " + err.Error())
	}
	user.Password = string(hashedPassword)

	// Set user defaults
	user.Type = "user"
	user.Created_at = time.Now().Format("2006-01-02 15:04:05")
	user.Updated_at = user.Created_at

	// Create new user
	err = repository.CreateUser(user)
	if err != nil {
		return errors.New("Error creating user: " + err.Error())
	}

	return nil
}
