package entity

import (
	"time"
)

type Todo struct {
	Entity

	Title      string     `json:"title" gorm:"size:128"`
	Content    string     `json:"content"`
	Importance bool       `json:"importance"`
	Deadline   *time.Time `json:"deadline"`
	Notified   bool       `json:"notified"`
	NotifyAt   *time.Time `json:"notifyAt"`
	Done       bool       `json:"done"`
	DoneAt     *time.Time `json:"doneAt"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoListID int64    `json:"todolistID"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`

	TodoRepeatPlanID int64          `json:"-"`
	TodoRepeatPlan   TodoRepeatPlan `json:"todoRepeatPlan"`

	TodoRepeatFromID *int64 `json:"todoRepeatFromID"` // last todo id
	TodoRepeatFrom   *Todo  `json:"-"`
}
