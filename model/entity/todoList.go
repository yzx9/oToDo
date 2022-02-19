package entity

type TodoList struct {
	Entity

	Name      string `json:"name" gorm:"size:128"`
	IsBasic   bool   `json:"-"`
	IsSharing bool   `json:"isSharing"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoListFolderID int64          `json:"todoListFolderID"`
	TodoListFolder   TodoListFolder `json:"-"`

	SharedUsers []*User `json:"-" gorm:"many2many:todo_list_shared_users"`
}
