package entity

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	Todos []Todo `json:"-"`
}
