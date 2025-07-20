package firestore_collections

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) CreateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Create Chapter: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) GetChapterById(novelId string, chapterId string, ctx context.Context) (*domain.Chapter, error) {
	doc, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapterId).Get(ctx)
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Chapter By Id: " + err.Error()), Status: http.StatusInternalServerError}
	}

	chapter := domain.Chapter{}
	if err = doc.DataTo(&chapter); err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Chapter By Id - Chapter Not Found"), Status: http.StatusNotFound}
		}
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get Chapter By Id: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return &chapter, nil
}

func (c *Client) GetAllChapters(novelId string, ctx context.Context) (*[]domain.Chapter, error) {
	doc, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Documents(ctx).GetAll()
	if err != nil {
		return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get All Chapters: " + err.Error()), Status: http.StatusInternalServerError}
	}

	chapters := []domain.Chapter{}
	for _, d := range doc {
		chapter := domain.Chapter{}
		if err = d.DataTo(&chapter); err != nil {
			return nil, &cmn.Error{Err: errors.New("Firestore Client Error - Get All Chapters: " + err.Error()), Status: http.StatusInternalServerError}
		}
		chapters = append(chapters, chapter)
	}

	return &chapters, nil
}

func (c *Client) UpdateChapter(novelId string, chapter domain.Chapter, ctx context.Context) error {
	updates := make(map[string]any)

	if chapter.Title != "" {
		updates["Title"] = chapter.Title
	}

	if chapter.Description != "" {
		updates["Description"] = chapter.Description
	}

	if chapter.Content != "" {
		updates["Content"] = chapter.Content
	}

	if len(updates) == 0 {
		return nil
	}

	updates["updatedAt"] = time.Now().Format("2006-01-02 15:04:05")

	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, chapter)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Update Chapter: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}

func (c *Client) DeleteChapter(novelId string, chapterId string, ctx context.Context) error {
	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapterId).Delete(ctx)
	if err != nil {
		return &cmn.Error{Err: errors.New("Firestore Client Error - Delete Chapter: " + err.Error()), Status: http.StatusInternalServerError}
	}

	return nil
}
