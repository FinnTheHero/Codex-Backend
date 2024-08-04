package aws

import (
	"Codex-Backend/api/utils"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateSession() (*dynamodb.DynamoDB, error) {

	APIKeys := utils.GetAPIKeys()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(APIKeys.Region),
		Credentials: credentials.NewStaticCredentials(
			APIKeys.AccessKey,
			APIKeys.SecretAccessKey,
			"",
		),
		Endpoint: aws.String("http://localhost:"),
	})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error creating session: %v", err))
	}

	svc := dynamodb.New(sess)

	return svc, nil
}
