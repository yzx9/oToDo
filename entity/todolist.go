package entity

import (
	"time"

	"github.com/google/uuid"
)

type TodoList struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string
	CreatedAt time.Time

	UserID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User
}
