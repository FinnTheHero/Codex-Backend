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
		return table.CreateTable("Users")
	}

	return nil
}

func GetUser(email string) (domain.User, error) {
	user := domain.User{}

	listOfUsers, err := getAllUsers()
	if err != nil {
		return domain.User{}, err
	}

	for _, u := range listOfUsers {
		if u.Email == email {
			user = u
			break
		}
	}

	return user, nil
}

func getAllUsers() ([]domain.User, error) {
	db, err := database.GetDynamoDBSession()

	if err != nil {
		return []domain.User{}, err
	}

	result, err := db.Scan(&dynamodb.ScanInput{
		TableName: aws.String("Users"),
	})

	if err != nil {
		return []domain.User{}, err
	}

	users := []domain.User{}

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		return []domain.User{}, err
	}

	return users, nil
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

	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(user.ID),
			},
			"Username": {
				S: aws.String(user.Username),
			},
			"Password": {
				S: aws.String(user.Password),
			},
			"Email": {
				S: aws.String(user.Email),
			},
			"Type": {
				S: aws.String(user.Type),
			},
			"CreatedAt": {
				S: aws.String(user.CreatedAt),
			},
			"UpdatedAt": {
				S: aws.String(user.UpdatedAt),
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
