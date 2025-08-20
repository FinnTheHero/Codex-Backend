package firestore_services

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure/collections"
	"context"
	"errors"
	"net/http"
	"time"
)

func GetCursorPaginatedChapters(options domain.CursorOptions, ctx context.Context) (*domain.CursorResponse, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	if options.Limit > 100 || options.Limit <= 0 {
		options.Limit = 100
	}

	response, err := c.CursorPagination(options, ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func BatchUploadChapters(novelId string, chapters []domain.Chapter, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	if len(chapters) == 0 {
		return &cmn.Error{
			Err:    errors.New("Nothing to upload"),
			Status: http.StatusInternalServerError,
		}
	}

	err = c.BatchUploadChapters(novelId, chapters, ctx)
	if err != nil {
		return err
	}

	return nil
}

func CreateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	id, err := cmn.GenerateID("chapter")
	if err != nil {
		return err
	}

	chapter.ID = id
	chapter.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	chapter.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	chapter.Deleted = false

	err = c.CreateChapter(novelId, chapter, ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetChapter(novelId, chapterId string, ctx context.Context) (*domain.Chapter, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	chapter, err := c.GetChapterById(novelId, chapterId, ctx)
	if err != nil {
		return nil, err
	}

	if chapter == nil {
		return nil, &cmn.Error{Err: errors.New("Chapter Service Error - Get Chapter - Chapter With ID " + chapterId + " In Novel With ID " + novelId + " Not Found"), Status: http.StatusNotFound}
	}

	return chapter, nil
}

func GetAllChapters(novelId string, ctx context.Context) (*[]domain.Chapter, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	chapters, err := c.GetAllChapters(novelId, ctx)
	if err != nil {
		return nil, err
	}

	if len(*chapters) == 0 {
		return nil, &cmn.Error{Err: errors.New("Chapter Service Error - Get All Chapters - Chapters In Novel With ID " + novelId + " Not Found"), Status: http.StatusNotFound}
	}

	return chapters, nil
}

func UpdateChapter(novelId string, chapter *domain.Chapter, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.UpdateChapter(novelId, *chapter, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteChapter(novelId, chapterId string, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.DeleteChapter(novelId, chapterId, ctx)
	if err != nil {
		return err
	}

	return nil
}
