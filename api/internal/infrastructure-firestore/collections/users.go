package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
)

func (c *Client) createUser(user domain.User, ctx context.Context) error {
	_, err := c.Client.Collection("Users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
