package entity

import "time"

type TodoStep struct {
	Entity

	Name   string    `json:"name" gorm:"-"`
	Done   bool      `json:"done"`
	DoneAt time.Time `json:"done_at"`

	TodoID string `json:"todo_id" gorm:"size:36"`
	Todo   Todo   `json:"-"`
}
