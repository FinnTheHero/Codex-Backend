package domain

// Novel struct used on backend
type Novel struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CreatedAt   string `json:"creation_date"`
	UpdatedAt   string `json:"update_date"`
	Index       int    `json:"index"`
	Deleted     bool   `json:"deleted"`
}

// Novel struct used on frontend
type FrontendNovel struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CreatedAt   string `json:"creation_date"`
	UpdatedAt   string `json:"update_date"`
}
