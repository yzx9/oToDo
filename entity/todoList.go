package entity

type TodoList struct {
	Entity

	Name      string `json:"name" gorm:"size:128"`
	Deletable bool   `json:"deletable"`

	UserID string `json:"userID" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TodoListFolderID string         `json:"todoListFolderID" gorm:"type:char(36);"`
	TodoListFolder   TodoListFolder `json:"-"`
}
