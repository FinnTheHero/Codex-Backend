package repository

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/database"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/* Find and return chapter from novel */
func GetChapter(novel, title string) (domain.Chapter, error) {

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return domain.Chapter{}, err
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(novel),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(title),
			},
		},
	})

	if err != nil {
		return domain.Chapter{}, err
	}

	if result.Item == nil {
		return domain.Chapter{}, errors.New("Could not find '" + title + "'")
	}

	chapter := domain.Chapter{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &chapter)
	if err != nil {
		return domain.Chapter{}, err
	}

	return chapter, nil
}

/* Return every chapter from novel */
func GetAllChapters(novel string) ([]domain.Chapter, error) {

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return nil, err
	}

	result, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(novel),
	})
	if err != nil {
		return nil, err
	}

	chapters := []domain.Chapter{}

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
func CreateChapter(novel string, chapter domain.Chapter) error {

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return err
	}

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

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
