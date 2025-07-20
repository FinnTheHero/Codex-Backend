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

func CreateNovel(novel domain.Novel, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	id, err := cmn.GenerateID("novel")
	if err != nil {
		return err
	}

	n, err := c.GetNovelById(id, ctx)
	if e, ok := err.(*cmn.Error); !ok {
		if e.StatusCode() != 404 {
			return err
		}
	}

	if n != nil {
		return &cmn.Error{Err: errors.New("Novel Service Error - Create Novel - Novel with ID " + id + " already exists"), Status: http.StatusConflict}
	}

	novel.ID = id
	novel.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	novel.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	novel.UploadedAt = time.Now().Format("2006-01-02 15:04:05")

	err = c.CreateNovel(novel, ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetNovel(id string, ctx context.Context) (*domain.Novel, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novel, err := c.GetNovelById(id, ctx)
	if err != nil {
		return nil, err
	}

	if novel == nil {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get Novel - Novel with ID " + id + " not found"), Status: http.StatusNotFound}
	}

	return novel, nil
}

func GetAllNovels(ctx context.Context) (*[]domain.Novel, error) {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novels, err := c.GetAllNovels(ctx)
	if err != nil {
		return nil, err
	}

	if len(*novels) == 0 {
		return nil, &cmn.Error{Err: errors.New("Novel Service Error - Get All Novels - No novels found"), Status: http.StatusNotFound}
	}

	return novels, nil
}

func UpdateNovel(id string, novel domain.Novel, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novel.ID = id
	novel.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	err = c.UpdateNovel(novel, ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteNovel(id string, ctx context.Context) error {
	client, err := firestore_client.FirestoreClient()
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	err = c.DeleteNovel(id, ctx)
	if err != nil {
		return err
	}

	return nil
}
