package entity

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	ID        string         `json:"id" gorm:"primaryKey;size:36"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
