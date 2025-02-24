package repository

import (
	"Codex-Backend/api/internal/domain"
	"Codex-Backend/api/internal/infrastructure/database"
	"Codex-Backend/api/internal/infrastructure/table"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func VerifyUsersTable() error {
	tableExists, err := table.IsTableCreated("Users")
	if err != nil {
		return err
	}

	if !tableExists {
		return table.CreateUsersTable()
	}

	return nil
}

/*
Find and return user.

Returns UserDTO as interface.
*/
func GetUser(email string) (any, error) {

	err := VerifyUsersTable()
	if err != nil {
		return nil, errors.New("Could not verify Users table" + err.Error())
	}

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return nil, err
	}

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Users"),
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("User not found")
	}

	user := domain.UserDTO{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

/* Create new user */
func CreateUser(user domain.User) error {

	err := VerifyUsersTable()
	if err != nil {
		return errors.New("Could not verify Users table" + err.Error())
	}

	db, err := database.GetDynamoDBSession()

	if err != nil {
		return err
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(user.Email),
			},
			"User": {
				M: av,
			},
		},
		TableName: aws.String("Users"),
	}

	_, err = db.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
