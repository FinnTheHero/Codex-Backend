package firestore_collections

import (
	cmn "Codex-Backend/api/internal/common"
	"Codex-Backend/api/internal/domain"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) BatchUploadChapters(novelId string, chapters []domain.Chapter, ctx context.Context) error {
	coll := c.Client.Collection("novels").Doc(novelId).Collection("chapters")
	const chunkSize = 500

	for i := 0; i < len(chapters); i += chunkSize {
		subset := chapters[i:min(i+chunkSize, len(chapters))]

		bw := c.Client.BulkWriter(ctx)
		jobs := make([]*firestore.BulkWriterJob, 0, len(subset))

		for _, chap := range subset {
			job, err := bw.Set(coll.Doc(chap.ID), chap)
			if err != nil {
				return &cmn.Error{
					Err:    fmt.Errorf("Firestore Client Error - Batch Upload Chapters - Enqueue failed for chapter %s: %w", chap.ID, err),
					Status: http.StatusInternalServerError,
				}
			}
			jobs = append(jobs, job)
		}

		bw.Flush()
		bw.End()

		// Check each job’s result to catch silent failures
		for j, job := range jobs {
			if _, err := job.Results(); err != nil {
				chap := subset[j]
				return &cmn.Error{
					Err:    fmt.Errorf("Firestore Client Error - Batch Upload Chapters - Write failed for chapter %s: %w", chap.ID, err),
					Status: http.StatusInternalServerError,
				}
			}
		}
	}

	return nil
}

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

	_, err := c.Client.Collection("novels").Doc(novelId).Collection("chapters").Doc(chapter.ID).Set(ctx, updates, firestore.MergeAll)
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
