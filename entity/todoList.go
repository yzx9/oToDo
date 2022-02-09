package entity

import (
	"time"
)

type TodoList struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Deletable bool      `json:"deletable"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	TodoListFolderID string         `json:"todo_list_folder_id"`
	TodoListFolder   TodoListFolder `json:"-"`
}
