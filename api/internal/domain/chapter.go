package domain

import "time"

type Chapter struct {
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	UploadDate   time.Time `json:"upload_date"`
	UpdateDate   time.Time `json:"update_date"`
	Content      string    `json:"content"`
}

type ChapterDTO struct {
	Title   string
	Author  string
	Chapter Chapter
}
