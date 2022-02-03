package entity

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []Tag     `json:"tags"`
	CreatedAt time.Time `json:"created_at"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`

	TodoListID uuid.UUID `json:"todolist_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	TodoList   TodoList  `json:"-"`
}
