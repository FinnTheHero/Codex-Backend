package repository

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/database"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetNovel(id string) (domain.Novel, error) {

	db, err := database.GetDynamoDBSession()
	if err != nil {
		return domain.Novel{}, err
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Novels"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		return domain.Novel{}, err
	}

	if result.Item == nil {
		return domain.Novel{}, errors.New("Could not find '" + id + "'")
	}

	novel := domain.Novel{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &novel)
	if err != nil {
		return novel, err
	}

	if novel.Title == "" {
		return domain.Novel{}, errors.New(id + " Not Found")
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

	novels := []domain.Novel{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &novels)
	if err != nil {
		return nil, err
	}

	if len(novels) == 0 {
		return nil, errors.New("No novels found")
	}

	return novels, nil
}

/* Add the new Novel to the Novels table */
func CreateNovel(novel domain.Novel) error {

	db, err := database.GetDynamoDBSession()
	if err != nil {
		return err
	}

	tableName := "Novels"

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(novel.ID),
			},
			"Title": {
				S: aws.String(novel.Title),
			},
			"Author": {
				S: aws.String(novel.Author),
			},
			"Description": {
				S: aws.String(novel.Description),
			},
			"CreationDate": {
				S: aws.String(novel.CreatedAt),
			},
			"UploadDate": {
				S: aws.String(novel.UploadedAt),
			},
			"UpdateDate": {
				S: aws.String(novel.UpdatedAt),
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
