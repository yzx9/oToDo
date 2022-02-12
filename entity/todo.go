package entity

import (
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Deadline  time.Time `json:"deadline"`
	Notified  bool      `json:"notified"`
	NotifyAt  time.Time `json:"notify_at"`
	Done      bool      `json:"done"`
	DoneAt    time.Time `json:"done_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID string `json:"user_id"`
	User   User   `json:"-"`

	TodoListID string   `json:"todolist_id"`
	TodoList   TodoList `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`

	RepeatPlanID string         `json:"-"`
	RepeatPlan   TodoRepeatPlan `json:"repeat_plan"`
	RepeatFrom   string         `json:"repeat_from"` // last todo id
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

type TodoRepeatPlanType string

const (
	TodoRepeatPlanTypeDay   TodoRepeatPlanType = "day"
	TodoRepeatPlanTypeWeek  TodoRepeatPlanType = "week"
	TodoRepeatPlanTypeMonth TodoRepeatPlanType = "month"
	TodoRepeatPlanTypeYear  TodoRepeatPlanType = "year"
)

type TodoRepeatPlan struct {
	ID       string    `json:"-"`
	Type     string    `json:"type"`
	Interval int       `json:"interval"`
	Before   time.Time `json:"before"`
	Weekday  [7]byte   `json:"weekday"` // Follow time.Weekend: Sunday Monday Tuesday Wednesday Thursday Friday Saturday
}
