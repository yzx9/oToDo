package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Name      string
	Password  []byte
	Email     string
	Avatar    string
	CreatedAt time.Time
}
