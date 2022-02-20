package entity

type Tag struct {
	Entity

	Name string `json:"name" gorm:"size:32;index:idx_tags_user,unique"`

	UserID int64 `json:"userID" gorm:"index:idx_tags_user,unique"`
	User   User  `json:"-"`

	Todos []Todo `json:"-" gorm:"many2many:tag_todos;"`
}
