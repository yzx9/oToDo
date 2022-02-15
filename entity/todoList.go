package entity

type TodoList struct {
	Entity

	Name      string `json:"name" gorm:"size:128"`
	Deletable bool   `json:"deletable"`

	UserID string `json:"user_id" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TodoListFolderID string         `json:"todo_list_folder_id" gorm:"type:char(36);"`
	TodoListFolder   TodoListFolder `json:"-"`
}
