package aws_methods

import (
	a "Codex-Backend/api/aws"
	"Codex-Backend/api/types"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetNovel(title string) (types.NovelSchema, error) {
	svc := a.Svc

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
		return res, err
	}

	if result.Item == nil {
		return res, errors.New("Could not find '" + title + "'")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &res)
	if err != nil {
		return res, err
	}

	if res.Title == "" {
		return res, errors.New(title + " Not Found")
	}

	return res, nil
}

func GetAllNovels() ([]types.NovelSchema, error) {
	svc := a.Svc

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Novels"),
	})

	if err != nil {
		return nil, err
	}

	novels := []types.NovelSchema{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &novels)
	if err != nil {
		return nil, err
	}

	if len(novels) == 0 {
		return nil, errors.New("No novels found")
	}

	return novels, nil
}

func GetChapter(novelTitle, chapterTitle string) (types.Chapter, error) {
	svc := a.Svc

	chapter := types.Chapter{}

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(novelTitle),
		Key: map[string]*dynamodb.AttributeValue{
			"Title": {
				S: aws.String(chapterTitle),
			},
		},
	})

	if err != nil {
		return chapter, err
	}

	if result.Item == nil {
		return chapter, errors.New("Could not find '" + chapterTitle + "'")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &chapter)
	if err != nil {
		return chapter, err
	}

	return chapter, nil
}

func GetAllChapters(novelTitle string) ([]types.Chapter, error) {
	svc := a.Svc

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String(novelTitle),
	})
	if err != nil {
		return nil, err
	}

	chapters := []types.Chapter{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &chapters)
	if err != nil {
		return nil, err
	}

	if len(chapters) == 0 {
		return nil, errors.New("No chapters found")
	}

	return chapters, nil
}

func GetTables() ([]string, error) {
	svc := a.Svc

	tableNames := []string{}

	input := &dynamodb.ListTablesInput{}

	for {
		// Get the list of tables
		result, err := svc.ListTables(input)
		if err != nil {
			return nil, err
		}

		for _, n := range result.TableNames {
			tableNames = append(tableNames, *n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return tableNames, nil
}
