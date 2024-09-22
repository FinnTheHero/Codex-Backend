package aws_services

import (
	"Codex-Backend/api/models"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/joho/godotenv"
)

var (
	once     sync.Once
	instance *AWSSession
)

type AWSSession struct {
	Svc *dynamodb.DynamoDB
}

/* Get existing aws session or create new one */
func GetAWSSession() *AWSSession {
	once.Do(func() {
		instance = NewAWSSession()
	})
	return instance
}

/* Create a new session to interact with aws services */
func NewAWSSession() *AWSSession {
	APIKeys := getAPIKeys()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(APIKeys.Region),
		Credentials: credentials.NewStaticCredentials(APIKeys.AccessKey, APIKeys.SecretAccessKey, ""),
	})

	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	return &AWSSession{
		Svc: dynamodb.New(sess),
	}
}

// Get API keys from '.env' file or environment variables
func getAPIKeys() models.APIKEYS {

	// For local development
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	// For Heroku deployment
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretAccessKey == "" || region == "" {
		log.Println("AWS environment variables not set")
	}

	log.Println("AWS env variables set")

	return models.APIKEYS{
		AccessKey:       accessKey,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}
}
