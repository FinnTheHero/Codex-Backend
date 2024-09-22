package aws_services

import (
	"Codex-Backend/api/models"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/* Find and return chapter from novel */
func GetChapter(novel, title string) (interface{}, error) {

	svc := GetAWSSession().Svc

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(novel),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(title),
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("Could not find '" + title + "'")
	}

	chapter := models.Chapter{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &chapter)
	if err != nil {
		return nil, err
	}

	return chapter, nil
}

/* Return every chapter from novel */
func GetAllChapters(novel string) (interface{}, error) {

	svc := GetAWSSession().Svc

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(novel),
	})
	if err != nil {
		return nil, err
	}

	chapters := []models.Chapter{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &chapters)
	if err != nil {
		return nil, err
	}

	if len(chapters) == 0 {
		return nil, errors.New("No chapters found for " + novel)
	}

	return chapters, nil
}

/* Add chapter to respective novel table. */
func CreateChapter(novel string, chapter models.Chapter) error {

	svc := GetAWSSession().Svc

	tableName := novel

	av, err := dynamodbattribute.MarshalMap(chapter)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(chapter.Title),
			},
			"Chapter": {
				M: av,
			},
		},
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
