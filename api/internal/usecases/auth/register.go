package auth_service

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/repository"
	"errors"
	"time"
)

func RegisterUser(credentials domain.NewUser) error {

	var newUser domain.User
	newUser.Email = credentials.Email
	newUser.Username = credentials.Username

	// Check if user already exists
	user, err := repository.GetUser(credentials.Email)
	if err != nil {
		if err.Error() != "User not found" {
			return err
		}
	}

	if user.Email == credentials.Email || user.Email != "" {
		return errors.New("Email already in use")
	}

	// Hash password
	hashedPassword, err := HashPassword(credentials.Password)
	if err != nil {
		return errors.New("Error hashing password: " + err.Error())
	}
	newUser.Password = string(hashedPassword)

	// Set user defaults
	newUser.Type = "user"
	newUser.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	newUser.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	newUser.ID, err = GenerateID("user")
	if err != nil {
		return errors.New("Error generating user ID: " + err.Error())
	}

	// Create new user
	err = repository.CreateUser(newUser)
	if err != nil {
		return errors.New("Error creating user: " + err.Error())
	}

	return nil
}
