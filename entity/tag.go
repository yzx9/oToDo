package entity

type Tag struct {
	Entity

	Name string `json:"name"`

	UserID string `json:"user_id" gorm:"size:36"`
	User   User   `json:"-"`

	Todos []Todo `json:"-" gorm:"many2many:tag_todos;"`
}
