package entity

import (
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	DoneAt    time.Time `json:"done_at"`
	Deadline  time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	IsRepeat       bool      `json:"is_repeat"`
	RepeatBefore   time.Time `json:"repeat_before"`
	RepeatInterval int       `json:"repeat_interval"` // 单位：秒
	RepeatFrom     string    `json:"repeat_from"`     // 上一次Todo ID

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
