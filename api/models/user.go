package models

import "github.com/golang-jwt/jwt/v5"

type User struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Type       string `json:"type"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDTO struct {
	Email string `json:"email"`
	User  User
}

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
