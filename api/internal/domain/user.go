package domain

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
