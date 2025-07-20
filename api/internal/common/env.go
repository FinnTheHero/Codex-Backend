package common

import (
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	return godotenv.Load(".env")
}

func GetEnvVariable(v string) (string, error) {
	env_variable := os.Getenv(v)
	if env_variable == "" {
		return "", &Error{Err: errors.New("Environmental Variable" + v + " Not Found"), Status: http.StatusNotFound}
	}

	return env_variable, nil
}
