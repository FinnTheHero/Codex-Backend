package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
)

func (c *Client) createChapterDocument(novel domain.Novel, chapter domain.Chapter, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novel.ID).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return err
	}

	return nil
}
