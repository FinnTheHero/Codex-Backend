package repository

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/database"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetNovel(title string) (any, error) {

	db, err := database.GetDynamoDBSession()
	if err != nil {
		return nil, err
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Novels"),
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

	novel := domain.NovelDTO{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &novel)
	if err != nil {
		return novel, err
	}

	if novel.Title == "" {
		return nil, errors.New(title + " Not Found")
	}

	return novel, nil
}

func GetAllNovels() (any, error) {

	db, err := database.GetDynamoDBSession()
	if err != nil {
		return nil, err
	}

	result, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Novels"),
	})
	if err != nil {
		return nil, err
	}

	novels := []domain.NovelDTO{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &novels)
	if err != nil {
		return nil, err
	}

	if len(novels) == 0 {
		return nil, errors.New("No novels found")
	}

	return novels, nil
}

/* Add Novel to the 'Novels' table */
func CreateNovel(novel domain.Novel) error {

	db, err := database.GetDynamoDBSession()
	if err != nil {
		return err
	}

	tableName := "Novels"

	av, err := dynamodbattribute.MarshalMap(novel)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(novel.Title),
			},
			"Author": {
				S: aws.String(novel.Author),
			},
			"Novel": {
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
