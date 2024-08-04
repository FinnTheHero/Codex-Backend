package utils

import (
	"Codex-Backend/api/types"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Get API keys from '.env' file or environment variables
func GetAPIKeys() types.AWSAPIKeys {
	// For local development
	err := godotenv.Load(".env")
	if err != nil {
		// Printf - as we dont want to crash app in heroku
		log.Println("Error loading .env file")
	}

	// For Heroku deployment
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccesstKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	output := os.Getenv("AWS_OUTPUT")

	if accessKey == "" || secretAccesstKey == "" || region == "" || output == "" {
		log.Println("AWS environment variables not set")
	}

	log.Println("AWS env variables set")

	return types.AWSAPIKeys{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccesstKey,
		Region:          region,
		Output:          output,
	}
}
