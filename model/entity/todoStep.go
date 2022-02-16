package entity

import "time"

type TodoStep struct {
	Entity

	Name   string     `json:"name" gorm:"-"`
	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	TodoID int64 `json:"todoID"`
	Todo   Todo  `json:"-"`
}
