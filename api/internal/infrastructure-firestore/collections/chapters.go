package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
)

func (c *Client) createChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return err
	}

	return nil
}
