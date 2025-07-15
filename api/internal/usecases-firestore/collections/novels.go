package firestore_services

import (
	"Codex-Backend/api/internal/domain"
	firestore_client "Codex-Backend/api/internal/infrastructure-firestore/client"
	firestore_collections "Codex-Backend/api/internal/infrastructure-firestore/collections"
	auth_service "Codex-Backend/api/internal/usecases/auth"
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateNovel(novel domain.Novel, ctx context.Context) error {
	client, err := firestore_client.NewFirestoreClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	id, err := auth_service.GenerateID("novel")
	if err != nil {
		return err
	}

	n, err := c.GetNovelById(id, ctx)
	if err != nil {
		return err
	}

	if n != nil {
		return status.Errorf(codes.AlreadyExists, "Novel with ID %s already exists", id)
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
	client, err := firestore_client.NewFirestoreClient(ctx)
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
		return nil, status.Errorf(codes.NotFound, "Novel with ID %s not found", id)
	}

	return novel, nil
}

func GetAllNovels(ctx context.Context) (*[]domain.Novel, error) {
	client, err := firestore_client.NewFirestoreClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	c := firestore_collections.Client{Client: client}

	novels, err := c.GetAllNovels(ctx)
	if err != nil {
		return nil, err
	}

	return novels, nil
}
