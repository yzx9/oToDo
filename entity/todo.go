package entity

import (
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []Tag     `json:"tags"`
	Done      bool      `json:"done"`
	DoneAt    time.Time `json:"done_at"`
	CreatedAt time.Time `json:"created_at"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	TodoListID string   `json:"todolist_id"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`
}

type TodoFile struct {
	ID       string `json:"id"`
	FileID   string `json:"file_id"`
	FileName string `json:"file_name" gorm:"-"`
	File     File   `json:"-"`

	TodoID string `json:"todo_id"`
	Todo   Todo   `json:"-"`
}

type TodoStep struct {
	ID     string    `json:"id"`
	Name   string    `json:"name" gorm:"-"`
	Done   bool      `json:"done"`
	DoneAt time.Time `json:"done_at"`

	TodoID string `json:"todo_id"`
	Todo   Todo   `json:"-"`
}
