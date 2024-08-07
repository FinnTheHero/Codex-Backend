package aws

import (
	"Codex-Backend/api/types"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
)

var Svc *dynamodb.DynamoDB

func init() {

	APIKeys := getAPIKeys()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(APIKeys.Region),
		Credentials: credentials.NewStaticCredentials(APIKeys.AccessKey, APIKeys.SecretAccessKey, ""),
	})

	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	Svc = dynamodb.New(sess)
}

// Get API keys from '.env' file or environment variables
func getAPIKeys() types.AWSAPIKeys {
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
