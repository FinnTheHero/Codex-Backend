package domain

import (
	"time"
)

type CursorOptions struct {
	NovelID   string `json:"novel_id"`
	Cursor    string `json:"cursor"`
	Limit     int    `json:"limit"`
	Ascending bool   `json:"sort_by"`
}

type CursorResponse struct {
	Chapters   []Chapter `json:"chapters"`
	NextCursor string    `json:"next_cursor"`
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
