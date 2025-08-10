package domain

import "cloud.google.com/go/firestore"

type CursorOptions struct {
	NovelID string              `json:"novel_id"`
	Cursor  string              `json:"cursor"`
	Limit   int                 `json:"limit"`
	SortBy  firestore.Direction `json:"sort_by"`
}

type CursorResponse struct {
	Chapters   []Chapter `json:"chapters"`
	NextCursor string    `json:"next_cursor"`
}

type Chapter struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CreatedAt   string `json:"creation_date"`
	UploadedAt  string `json:"upload_date"`
	UpdatedAt   string `json:"update_date"`
	Content     string `json:"content"`
}
