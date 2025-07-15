package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
)

func (c *Client) CreateUser(user domain.User, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
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

	return nil, errors.New("User not found")
}

func (c *Client) GetUserById(userId string, ctx context.Context) (*domain.User, error) {
	doc, err := c.Client.Collection("users").Doc(userId).Get(ctx)
	if err != nil {
		return nil, err
	}

	user := domain.User{}
	if err = doc.DataTo(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) GetAllUsers(ctx context.Context) (*[]domain.User, error) {
	doc, err := c.Client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	users := []domain.User{}
	for _, d := range doc {
		var user domain.User
		err = d.DataTo(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func (c *Client) UpdateUser(user domain.User, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteUser(id string, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
