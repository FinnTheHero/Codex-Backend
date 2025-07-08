package table

import (
	"Codex-Backend/api/internal/infrastructure/database"
	"slices"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func IsTableCreated(id string) (bool, error) {
	ids, err := GetTableIds()
	if err != nil {
		return false, err
	}

	if slices.Contains(ids, id) {
		return true, nil
	}

	return false, nil
}

func GetTableIds() ([]string, error) {
	db, err := database.GetDynamoDBSession()
	if err != nil {
		return nil, err
	}

	ids := []string{}

	input := &dynamodb.ListTablesInput{}

	for {
		// Get the list of tables
		result, err := db.ListTables(input)
		if err != nil {
			return nil, err
		}

		for _, n := range result.TableNames {
			ids = append(ids, *n)
		}

		// assign the last read tablename as the start for our next call to the ListTables function
		// the maximum number of table names returned in a call is 100 (default), which requires us to make
		// multiple calls to the ListTables function to retrieve all table names
		input.ExclusiveStartTableName = result.LastEvaluatedTableName

		if result.LastEvaluatedTableName == nil {
			break
		}
	}

	return ids, nil
}

func CreateTable(id string) error {
	db, err := database.GetDynamoDBSession()
	if err != nil {
		return err
	}

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("ID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("ID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(id),
	}

	_, err = db.CreateTable(input)
	return err
}

func DeleteTable(id string) error {
	db, err := database.GetDynamoDBSession()
	if err != nil {
		return err
	}

	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(id),
	}

	_, err = db.DeleteTable(input)
	if err != nil {
		return err
	}

	return nil
}
