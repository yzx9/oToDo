package entity

type Tag struct {
	Entity

	Name string `json:"name" gorm:"size:32"`

	UserID string `json:"user_id" gorm:"type:char(36);"`
	User   User   `json:"-"`

	Todos []Todo `json:"-" gorm:"many2many:tag_todos;"`
}
