package service

import (
	cmn "Codex-Backend/api/common"
	db "Codex-Backend/api/internal/database"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
)

func GetPaginatedChapters(options domain.CursorOptions, ctx context.Context) (*domain.CursorResponse, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	chapters, nextCursor, err := client.ListChaptersSeek(options, ctx)
	if err != nil {
		return nil, err
	}

	response := &domain.CursorResponse{
		Chapters:   chapters,
		NextCursor: nextCursor,
	}

	return response, nil
}

func CreateChapter(chapter domain.CreateChapter, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.CreateChapter(chapter, ctx); err != nil {
		return err
	}

	return nil
}

func GetChapter(novelId, chapterId string, ctx context.Context) (domain.Chapter, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return domain.Chapter{}, err
	}

	chapter, err := client.GetChapterById(novelId, chapterId, ctx)
	if err != nil {
		return domain.Chapter{}, err
	}

	return chapter, nil
}

func GetAllChapters(novelId string, pageSize int, asc bool, ctx context.Context) ([]domain.Chapter, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	chapters, err := client.GetAllChapters(novelId, pageSize, asc, ctx)
	if err != nil {
		return nil, err
	}

	if len(chapters) == 0 {
		return nil, &cmn.Error{Err: errors.New("Chapter Service Error - Get All Chapters - Chapters In Novel With ID " + novelId + " Not Found"), Status: http.StatusNotFound}
	}

	return chapters, nil
}

func UpdateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.UpdateChapter(novelId, chapter, ctx); err != nil {
		return err
	}

	return nil
}

func DeleteChapter(novelId, chapterId string, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.DeleteChapter(novelId, chapterId, ctx); err != nil {
		return err
	}

	return nil
}
