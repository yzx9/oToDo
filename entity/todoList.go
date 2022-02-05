package entity

import (
	"time"

	"github.com/google/uuid"
)

type TodoList struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`
}
