package firestore_collections

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
)

func (c *Client) CreateUser(user domain.User, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Creating User: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) GetUserByEmail(email string, ctx context.Context) (*domain.User, error) {
	users, err := c.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, user := range *users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, &cmn.Error{Err: errors.New("User not found"), Status: http.StatusNotFound}
}

func (c *Client) GetUserById(userId string, ctx context.Context) (*domain.User, error) {
	doc, err := c.Client.Collection("users").Doc(userId).Get(ctx)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Getting User by ID: " + err.Error()), Status: http.StatusInternalServerError}
	}

	user := domain.User{}
	if err = doc.DataTo(&user); err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Getting User by ID: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return &user, nil
}

func (c *Client) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	doc, err := c.Client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Getting All Users: " + err.Error()), Status: http.StatusInternalServerError}
	}

	users := []domain.User{}
	for _, d := range doc {
		var user domain.User
		err = d.DataTo(&user)
		if err != nil {
			return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Getting All Users: " + err.Error()), Status: http.StatusInternalServerError}
		}
		users = append(users, user)
	}

	return &users, nil
}

func (c *Client) UpdateUser(user domain.User, ctx context.Context) error {
	updates := make(map[string]any)

	if user.Email != "" {
		updates["Email"] = user.Email
	}

	if user.Username != "" {
		updates["Username"] = user.Username
	}

	if user.Password != "" {
		updates["Password"] = user.Password
	}

	if user.Type != "" {
		updates["Type"] = user.Type
	}

	if len(updates) == 0 {
		return nil
	}

	updates["UpdatedAt"] = user.UpdatedAt

	_, err := c.Client.Collection("users").Doc(user.ID).Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Updating User: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) DeleteUser(id string, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(id).Delete(ctx)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Deleting User: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}
