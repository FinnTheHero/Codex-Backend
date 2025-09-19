package domain

import (
	"time"

	"cloud.google.com/go/firestore"
)

type CursorOptions struct {
	NovelID string              `json:"novel_id"`
	Cursor  int                 `json:"cursor"`
	Limit   int                 `json:"limit"`
	SortBy  firestore.Direction `json:"sort_by"`
}

type CursorResponse struct {
	Chapters   []FrontendChapter `json:"chapters"`
	NextCursor int               `json:"next_cursor"`
}

// Chapter struct used on backend
type Chapter struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"creation_date"`
	UpdatedAt   time.Time `json:"update_date"`
	Content     string    `json:"content"`
	Index       int       `json:"index"`
	Deleted     bool      `json:"deleted"`
}

// Chapter struct used on frontend
type FrontendChapter struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"update_date"`
	Content   string    `json:"content"`
}
