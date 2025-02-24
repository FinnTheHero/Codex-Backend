package error_service

import (
	"errors"
)


var (

	ErrEmailTaken = errors.New("Email already in use")
	ErrUserNotFound = errors.New("User not found")

)
