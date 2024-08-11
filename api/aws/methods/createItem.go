package aws_methods

import (
	a "Codex-Backend/api/aws"
	"Codex-Backend/api/types"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Creates a table with the title as the table name.
//
// `title` can only contain letters, numbers, underscores, dot and hyphens.
func CreateTable(title string) error {
	svc := a.Svc
	tableName := title

	finalTableName := strings.ReplaceAll(tableName, " ", "_")

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Title"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Title"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(finalTableName),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		return err
	}

	return nil
}

// Adds `novel` to the 'Novels' table.
func CreateNovel(novel types.Novel) error {
	svc := a.Svc
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

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

// Adds chapter to respective novel table.
//
// `novelTitle` is the title of the novel to which the `chapter` is to be added.
func CreateChapter(novelTitle string, chapter types.Chapter) error {
	svc := a.Svc
	tableName := novelTitle

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
