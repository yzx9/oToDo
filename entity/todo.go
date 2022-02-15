package entity

import (
	"time"
)

type Todo struct {
	Entity

	Title      string    `json:"title" gorm:"size:128"`
	Content    string    `json:"content"`
	Importance bool      `json:"importance"`
	Deadline   time.Time `json:"deadline"`
	Notified   bool      `json:"notified"`
	NotifyAt   time.Time `json:"notify_at"`
	Done       bool      `json:"done"`
	DoneAt     time.Time `json:"done_at"`

	UserID string `json:"user_id" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TodoListID string   `json:"todolist_id" gorm:"type:char(36);"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`

	TodoRepeatPlanID string         `json:"-" gorm:"type:char(36);"`
	TodoRepeatPlan   TodoRepeatPlan `json:"todo_repeat_plan"`

	TodoRepeatFromID *string `json:"todo_repeat_from_id" gorm:"type:char(36);"` // last todo id
	TodoRepeatFrom   *Todo   `json:"-"`
}
