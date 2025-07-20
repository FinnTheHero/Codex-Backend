package common

import (
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, &Error{Err: errors.New("Password Service Error - Hash Password: " + err.Error()), Status: http.StatusInternalServerError}
	}
	return hash, nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return &Error{Err: errors.New("Password Service Error - Verify Password: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}
