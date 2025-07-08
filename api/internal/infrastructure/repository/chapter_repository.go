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
func GetAllChapters(novelId string) ([]domain.Chapter, error) {

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return nil, err
	}

	result, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String(novelId),
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
		return nil, errors.New("No chapters found for " + novelId)
	}

	return chapters, nil
}

func CreateChapter(novel string, chapter domain.Chapter) error {
	db, err := database.GetDynamoDBSession()

	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(chapter.ID),
			},
			"Title": {
				S: aws.String(chapter.Title),
			},
			"Author": {
				S: aws.String(chapter.Author),
			},
			"Description": {
				S: aws.String(chapter.Description),
			},
			"CreationDate": {
				S: aws.String(chapter.CreatedAt),
			},
			"UploadDate": {
				S: aws.String(chapter.UploadedAt),
			},
			"UpdateDate": {
				S: aws.String(chapter.UpdatedAt),
			},
			"Content": {
				S: aws.String(chapter.Content),
			},
		},
		TableName: aws.String(novel),
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
