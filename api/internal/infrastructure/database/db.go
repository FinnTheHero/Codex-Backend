package database

import (
	aws_session "Codex-Backend/api/internal/infrastructure/aws"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	once sync.Once
	instance *dynamodb.DynamoDB
	initErr error
)

func GetDynamoDBSession() (*dynamodb.DynamoDB, error) {
	once.Do(func() {
		instance, initErr = NewDynamoDBSession()
		if initErr != nil {
			log.Fatalf("Error creating aws session: %v", initErr)
		}
	})

	return instance, initErr
}

func NewDynamoDBSession() (*dynamodb.DynamoDB, error) {

	awsSession, err := aws_session.GetAWSSession()

	return dynamodb.New(awsSession), err
}
