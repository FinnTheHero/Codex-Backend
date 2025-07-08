package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AWSKeys struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`
}

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

func GetAWSKeys() (AWSKeys, error) {
	accessKey, err := GetEnvVariable("AWS_ACCESS_KEY")
	if err != nil {
		return AWSKeys{}, fmt.Errorf("Error getting %s: %v", accessKey, err)
	}
	secretKey, err := GetEnvVariable("AWS_SECRET_KEY")
	if err != nil {
		return AWSKeys{}, fmt.Errorf("Error getting %s: %v", secretKey, err)
	}
	region, err := GetEnvVariable("AWS_REGION")
	if err != nil {
		return AWSKeys{}, fmt.Errorf("Error getting %s: %v", region, err)
	}

	log.Println("AWS Env Variables Set")

	return AWSKeys{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}, nil
}
