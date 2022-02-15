package entity

import "time"

type TodoStep struct {
	Entity

	Name   string     `json:"name" gorm:"-"`
	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	TodoID string `json:"todoID" gorm:"type:char(36);"`
	Todo   Todo   `json:"-"`
}
