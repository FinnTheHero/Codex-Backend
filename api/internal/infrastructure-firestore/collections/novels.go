package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"

	"cloud.google.com/go/firestore"
)

type Client struct {
	*firestore.Client
}

func (c *Client) createNovelDocument(novel domain.Novel, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novel.ID).Set(ctx, novel)
	if err != nil {
		return err
	}

	return nil
}
