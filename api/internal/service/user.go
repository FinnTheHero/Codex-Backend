package service

import (
	cmn "Codex-Backend/api/common"
	firestore_client "Codex-Backend/api/internal/database/client"
	firestore_collections "Codex-Backend/api/internal/database/collections"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
	"time"
)

func LoginUser(credentials domain.Credentials, ctx context.Context) (*domain.User, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	user, err := c.GetUserByEmail(credentials.Email, ctx)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &cmn.Error{Err: errors.New("Login Service Error - User not found"), Status: http.StatusNotFound}
	}

	err = cmn.VerifyPassword(user.Password, credentials.Password)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Login Service Error - Invalid password"), Status: http.StatusUnauthorized}
	}

	return user, nil
}

func GetUserByID(id string, ctx context.Context) (*domain.User, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	user, err := c.GetUserById(id, ctx)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &cmn.Error{Err: errors.New("Get User By ID Service Error - User not found"), Status: http.StatusNotFound}
	}

	return user, nil
}

func RegisterUser(newUser domain.NewUser, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	user, err := c.GetUserByEmail(newUser.Email, ctx)
	if e, ok := err.(*cmn.Error); ok {
		if e.StatusCode() != http.StatusNotFound {
			return &cmn.Error{Err: errors.New("Register Service Error - Getting User By Email: " + err.Error()), Status: http.StatusInternalServerError}
		}
	}

	if user != nil {
		return &cmn.Error{Err: errors.New("Register Service Error - User With Email " + newUser.Email + " Already Exists"), Status: http.StatusConflict}
	}

	id, err := cmn.GenerateID("user")
	if err != nil {
		return err
	}

	hashedPassword, err := cmn.HashPassword(newUser.Password)
	if err != nil {
		return err
	}

	err = c.CreateUser(domain.User{
		ID:        id,
		Username:  newUser.Username,
		Password:  string(hashedPassword),
		Email:     newUser.Email,
		Type:      "User",
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}, ctx)
	if err != nil {
		return err
	}

	return nil
}

func LogoutUser(tokenString string) error {
	if tokenString == "" {
		return &cmn.Error{Err: errors.New("Logout Service Error - Token not found"), Status: http.StatusBadRequest}
	}

	return nil
}

func UpdateUser(updatedUser domain.User, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.UpdateUser(updatedUser, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id string, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.DeleteUser(id, ctx)
	if err != nil {
		return err
	}

	return nil
}
