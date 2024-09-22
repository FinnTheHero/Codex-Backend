package aws_services

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/* Checks if table exists */
func IsTableCreated(tableName string) (bool, error) {
	tableNames, err := GetTables()
	if err != nil {
		return false, err
	}

	for _, name := range tableNames {
		if name == tableName {
			return true, nil
		}
	}

	return false, nil
}

/* Returns a list of tables */
func GetTables() ([]string, error) {

	svc := GetAWSSession().Svc

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

/*
Creates a table with the title as the table name.

`title` can only contain letters, numbers, underscores, dot and hyphens  -  No spaces!
*/
func CreateTable(title string) error {

	svc := GetAWSSession().Svc

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

/* Create 'Users' table */
func CreateUsersTable() error {

	svc := GetAWSSession().Svc

	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String("Users"),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		return err
	}

	return nil
}
