package auth_services

import gonanoid "github.com/matoous/go-nanoid/v2"

func GenerateID() (string, error) {
	return gonanoid.Generate("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 6)
}
