package common

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(&Error{Err: errors.New("Failed to load environment variables"), Status: http.StatusInternalServerError})
	}
}

func GetEnvVariable(v string) string {
	env_variable := os.Getenv(v)
	if env_variable == "" {
		log.Fatal(&Error{Err: errors.New("Environmental Variable " + v + " Not Found"), Status: http.StatusNotFound})
	}

	return env_variable
}

func GetDomains(v string) []string {
	env_variable := os.Getenv(v)
	if env_variable == "" {
		log.Fatal(&Error{Err: errors.New("Environmental Variable " + v + " Not Found"), Status: http.StatusNotFound})
	}

	result := strings.Split(env_variable, ",")

	return result
}
