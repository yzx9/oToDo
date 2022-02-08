package entity

import (
	"time"

	"github.com/google/uuid"
)

type TodoListFolder struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`
}
