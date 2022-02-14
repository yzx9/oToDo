package entity

import (
	"time"
)

type Todo struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Importance bool      `json:"importance"`
	Deadline   time.Time `json:"deadline"`
	Notified   bool      `json:"notified"`
	NotifyAt   time.Time `json:"notify_at"`
	Done       bool      `json:"done"`
	DoneAt     time.Time `json:"done_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	TodoListID string   `json:"todolist_id"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`

	TodoRepeatPlanID string         `json:"-"`
	TodoRepeatPlan   TodoRepeatPlan `json:"todo_repeat_plan"`
	TodoRepeatFrom   string         `json:"todo_repeat_from"` // last todo id
}
