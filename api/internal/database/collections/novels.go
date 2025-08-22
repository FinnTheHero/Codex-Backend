package firestore_collections

import (
	cmn "Codex-Backend/api/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	*firestore.Client
}

func (c *Client) CreateNovel(novel domain.Novel, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novel.ID).Set(ctx, novel)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Create Novel: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) GetNovelById(id string, ctx context.Context) (*domain.Novel, error) {
	doc, err := c.Client.Collection("novels").Doc(id).Get(ctx)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Novel by ID: " + err.Error()), Status: http.StatusInternalServerError}
	}

	novel := domain.Novel{}
	if err := doc.DataTo(&novel); err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Novel by ID - Novel not found"), Status: http.StatusNotFound}
		}
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Novel by ID: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return &novel, nil
}

func (c *Client) GetAllNovels(ctx context.Context) (*[]domain.Novel, error) {
	doc, err := c.Client.Collection("novels").Documents(ctx).GetAll()
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get All Novels: " + err.Error()), Status: http.StatusInternalServerError}
	}

	novels := []domain.Novel{}
	for _, d := range doc {
		novel := domain.Novel{}
		if err := d.DataTo(&novel); err != nil {
			return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get All Novels: " + err.Error()), Status: http.StatusInternalServerError}
		}
		novels = append(novels, novel)
	}

	return &novels, nil
}

func (c *Client) UpdateNovel(novel domain.Novel, ctx context.Context) error {
	updates := make(map[string]any)

	if novel.Title != "" {
		updates["Title"] = novel.Title
	}

	if novel.Description != "" {
		updates["Description"] = novel.Description
	}

	if len(updates) == 0 {
		return nil
	}

	updates["UpdatedAt"] = time.Now().Format("2006-01-02 15:04:05")

	_, err := c.Client.Collection("novels").Doc(novel.ID).Set(ctx, updates, firestore.MergeAll)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Update Novel: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) DeleteNovel(novelId string, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Delete(ctx)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Delete Novel: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) GetNovelByTitle(title string, ctx context.Context) (*domain.Novel, error) {
	query := c.Client.Collection("novels").Where("Title", "==", title).Limit(1)
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Novel by Title: " + err.Error()), Status: http.StatusInternalServerError}
	}

	if len(docs) == 0 {
		return nil, &cmn.Error{Err: errors.New("Novel not found"), Status: http.StatusNotFound}
	}

	novel := domain.Novel{}
	if err := docs[0].DataTo(&novel); err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Novel by Title: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return &novel, nil
}
