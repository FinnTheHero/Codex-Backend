package common

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	return godotenv.Load(".env")
}

func GetEnvVariable(v string) (string, error) {
	env_variable := os.Getenv(v)
	if env_variable == "" {
		return "", fmt.Errorf("Environmental Variable \"%v\" Not Found", v)
	}

	return env_variable, nil
}
