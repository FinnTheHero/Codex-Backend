package aws_methods

import (
	"Codex-Backend/api/types"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetNovel(svc *dynamodb.DynamoDB, title string) (types.NovelSchema, error) {
	var res types.NovelSchema

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Novels"),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(title),
			},
		},
	})

	if err != nil {
		return res, errors.New(fmt.Sprintf("Got error calling GetItem: %s", err))
	}

	if result.Item == nil {
		return res, errors.New("Could not find '" + title + "'")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &res)
	if err != nil {
		return res, errors.New(fmt.Sprintf("Got error unmarshalling: %s", err))
	}

	return res, nil
}

func GetChapter(svc *dynamodb.DynamoDB, title string) (types.Chapter, error) {
	chapter := types.Chapter{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Chapters"),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(title),
			},
		},
	})

	if err != nil {
		return chapter, errors.New(fmt.Sprintf("Got error calling GetItem: %s", err))
	}

	if result.Item == nil {
		return chapter, errors.New("Could not find '" + title + "'")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &chapter)
	if err != nil {
		return chapter, errors.New(fmt.Sprintf("Got error unmarshalling: %s", err))
	}

	return chapter, nil
}
