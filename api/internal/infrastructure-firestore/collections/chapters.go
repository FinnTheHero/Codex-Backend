package firestore_collections

import (
	"Codex-Backend/api/internal/domain"
	"context"
)

func (c *Client) createChapterInNovel(novelId string, chapter domain.Chapter, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getChapterByIdFromNovel(novelId string, chapterId string, ctx context.Context) (domain.Chapter, error) {
	doc, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapterId).Get(ctx)
	if err != nil {
		return domain.Chapter{}, err
	}

	chapter := domain.Chapter{}
	if err = doc.DataTo(&chapter); err != nil {
		return domain.Chapter{}, err
	}

	return chapter, nil
}

func (c *Client) getAllChaptersFromNovel(novelId string, ctx context.Context) (*[]domain.Chapter, error) {
	doc, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	chapters := []domain.Chapter{}
	for _, d := range doc {
		chapter := domain.Chapter{}
		if err = d.DataTo(&chapter); err != nil {
			return nil, err
		}
		chapters = append(chapters, chapter)
	}

	return &chapters, nil
}

func (c *Client) updateChapterFromNovel(novelId string, chapter domain.Chapter, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteChapterFromNovel(novelId string, chapterId string, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapterId).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}
