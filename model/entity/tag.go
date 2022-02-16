package entity

type Tag struct {
	Entity

	Name string `json:"name" gorm:"size:32"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	Todos []Todo `json:"-" gorm:"many2many:tag_todos;"`
}
