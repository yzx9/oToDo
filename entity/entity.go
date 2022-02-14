package entity

import "time"

type Entity struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Deleted   bool      `json:"-"`
	DeletedAt time.Time `json:"-"`
}
