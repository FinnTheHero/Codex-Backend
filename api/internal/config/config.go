package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)


type AWSKeys struct {
	AccessKey       string `json:"access_key"`
	SecretKey       string `json:"secret_key"`
	Region          string `json:"region"`
}

func GetAWSKeys() (AWSKeys, error) {
	accessKey, err := GetEnvVariable("AWS_ACCESS_KEY")
	if err != nil {
		return AWSKeys{}, err
	}
	secretKey, err := GetEnvVariable("AWS_SECRET_KEY")
	if err != nil {
		return AWSKeys{}, err
	}
	region, err := GetEnvVariable("AWS_REGION")
	if err != nil {
		return AWSKeys{}, err
	}

	log.Println("AWS Env Variables Set")

	return AWSKeys{
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		Region:          region,
	}, nil
}

func LoadEnvVariables() error {
	return godotenv.Load(".env")
}

func GetEnvVariable(v string) (string, error) {
	env_variable := os.Getenv(v)
	if env_variable == "" {
		return "", errors.New(fmt.Sprintf("Environmental Variable \"%v\" Not Found", v))
	}

	return env_variable, nil
}
