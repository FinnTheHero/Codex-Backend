package firestore_services

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	auth_service "Codex-Backend/api/internal/usecases/auth"
	"context"
	"errors"
	"net/http"
	"time"
)

func CreateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	id, err := auth_service.GenerateID("chapter")
	if err != nil {
		return err
	}

	ch, err := c.GetChapterById(novelId, id, ctx)
	if err != nil {
		return err
	}

	if ch != nil {
		return &cmn.Error{Err: errors.New("Novel Service Error - Create Chapter - Chapter With ID " + id + " In Novel With ID " + novelId + " Already Exists"), Status: http.StatusConflict}
	}

	chapter.ID = id
	chapter.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	chapter.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	chapter.UploadedAt = time.Now().Format("2006-01-02 15:04:05")

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

	if chapters == nil {
		return nil, &cmn.Error{Err: errors.New("Chapter Service Error - Get All Chapters - Chapters In Novel With ID " + novelId + " Not Found"), Status: http.StatusNotFound}
	}

	return chapters, nil
}
