package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"

	"cloud.google.com/go/firestore"
)

type Client struct {
	*firestore.Client
}

func (c *Client) createNovel(novel domain.Novel, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novel.ID).Set(ctx, novel)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getNovelById(id string, ctx context.Context) (domain.Novel, error) {
	doc, err := c.Client.Collection("novels").Doc(id).Get(ctx)
	if err != nil {
		return domain.Novel{}, err
	}

	novel := domain.Novel{}
	if err := doc.DataTo(&novel); err != nil {
		return domain.Novel{}, err
	}

	return novel, nil
}

func (c *Client) getAllNovels(ctx context.Context) (*[]domain.Novel, error) {
	doc, err := c.Client.Collection("novels").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	novels := []domain.Novel{}
	for _, d := range doc {
		novel := domain.Novel{}
		if err := d.DataTo(&novel); err != nil {
			return nil, err
		}
		novels = append(novels, novel)
	}

	return &novels, nil
}
