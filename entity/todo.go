package entity

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title     string
	Content   string
	Tags      []Tag
	CreatedAt time.Time

	UserID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`

	TodoListID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	TodoList   TodoList  `json:"-"`
}
