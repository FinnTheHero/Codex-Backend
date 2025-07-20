package domain

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
