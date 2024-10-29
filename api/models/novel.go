package models

import "time"

type Novel struct {
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	UploadDate   time.Time `json:"upload_date"`
	UpdateDate   time.Time `json:"update_date"`
}

type NovelDTO struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Novel  Novel  `json:"novel"`
}
