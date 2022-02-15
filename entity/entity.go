package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID        string         `json:"id" gorm:"primaryKey;type:char(36);"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.NewString()
	return
}
