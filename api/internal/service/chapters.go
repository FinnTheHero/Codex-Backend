package service

import (
	cmn "Codex-Backend/api/common"
	db "Codex-Backend/api/internal/database"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
)

func GetCursorPaginatedChapters(options domain.CursorOptions, ctx context.Context) (*domain.CursorResponse, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	if options.Limit > 100 || options.Limit <= 0 {
		options.Limit = 100
	}

	response, err := client.CursorPagination(options, ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func BatchUploadChapters(novelId string, chapters []domain.Chapter, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if len(chapters) == 0 {
		return &cmn.Error{
			Err:    errors.New("Nothing to upload"),
			Status: http.StatusInternalServerError,
		}
	}

	if err = client.BatchUploadChapters(novelId, chapters, ctx); err != nil {
		return err
	}

	return nil
}

func CreateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	client, err := db.GetClient(ctx)
	if err != nil {
		return err
	}

	if err = client.CreateChapter(novelId, chapter, ctx); err != nil {
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

	if err = client.DeleteChapter(chapterId, ctx); err != nil {
		return err
	}

	return nil
}
