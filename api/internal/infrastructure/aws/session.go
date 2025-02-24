package aws_session

import (
	"Codex-Backend/api/internal/config"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

var (
	once     sync.Once
	instance *session.Session
	initErr  error
)

/* Get existing aws session or create new one */
func GetAWSSession() (*session.Session, error) {
	once.Do(func() {
		instance, initErr = NewAWSSession()
		if initErr != nil {
			log.Fatalf("Error creating aws session: %v", initErr)
		}
	})
	return instance, initErr
}

/* Create a new session to interact with aws services */
func NewAWSSession() (*session.Session, error) {
	AWSKeys, err := config.GetAWSKeys()
	if err != nil {
		log.Fatalf("Error getting AWS keys: %v", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWSKeys.Region),
		Credentials: credentials.NewStaticCredentials(AWSKeys.AccessKey, AWSKeys.SecretKey, ""),
	})

	return sess, err
}
