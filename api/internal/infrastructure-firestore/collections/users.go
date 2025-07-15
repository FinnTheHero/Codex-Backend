package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
)

func (c *Client) createUser(user domain.User, ctx context.Context) error {
	_, err := c.Client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getUserByEmail(email string, ctx context.Context) (*domain.User, error) {
	doc, err := c.Client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var user domain.User
	for _, d := range doc {
		if d.Data()["email"] == email {
			err = d.DataTo(&user)
			if err != nil {
				return nil, err
			}
			return &user, nil
		}
	}

	return nil, errors.New("User not found")
}

func (c *Client) getUserById(id string, ctx context.Context) (*domain.User, error) {
	doc, err := c.Client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c *Client) getAllUsers(ctx context.Context) (*[]domain.User, error) {
	doc, err := c.Client.Collection("users").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var users []domain.User
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
