package types

import "time"

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

type Novel struct {
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	UploadDate   time.Time `json:"upload_date"`
	UpdateDate   time.Time `json:"update_date"`
}

type Chapter struct {
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	CreationDate time.Time `json:"creation_date"`
	UploadDate   time.Time `json:"upload_date"`
	UpdateDate   time.Time `json:"update_date"`
	Content      string    `json:"content"`
}

type NovelSchema struct {
	Title  string
	Author string
	Novel  Novel
}

type ChapterSchema struct {
	Title   string
	Author  string
	Chapter Chapter
}

type AWSAPIKeys struct {
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_key"`
	Region          string `json:"region"`
	Output          string `json:"output"`
}
