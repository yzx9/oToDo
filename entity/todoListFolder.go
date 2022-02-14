package entity

import (
	"time"
)

type TodoListFolder struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	TodoLists []TodoList `json:"-"`
}
