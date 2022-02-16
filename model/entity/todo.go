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

	UserID string `json:"userID" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TodoListID string   `json:"todolistID" gorm:"type:char(36);"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`

	TodoRepeatPlanID string         `json:"-" gorm:"type:char(36);"`
	TodoRepeatPlan   TodoRepeatPlan `json:"todoRepeatPlan"`

	TodoRepeatFromID *string `json:"todoRepeatFromID" gorm:"type:char(36);"` // last todo id
	TodoRepeatFrom   *Todo   `json:"-"`
}
