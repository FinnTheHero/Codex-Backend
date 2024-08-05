package aws

import (
	"Codex-Backend/api/utils"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var Svc *dynamodb.DynamoDB

func init() {

	APIKeys := utils.GetAPIKeys()

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(APIKeys.Region),
		Credentials: credentials.NewStaticCredentials(APIKeys.AccessKey, APIKeys.SecretAccessKey, ""),
	})

	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}

	Svc = dynamodb.New(sess)
}
