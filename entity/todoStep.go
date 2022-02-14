package entity

import "time"

type TodoStep struct {
	ID     string    `json:"id"`
	Name   string    `json:"name" gorm:"-"`
	Done   bool      `json:"done"`
	DoneAt time.Time `json:"done_at"`

	TodoID string `json:"todo_id"`
	Todo   Todo   `json:"-"`
}
