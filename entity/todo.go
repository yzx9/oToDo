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
	Done      bool      `json:"done"`
	DoneAt    time.Time `json:"done_at"`
	CreatedAt time.Time `json:"created_at"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	User   User      `json:"-"`

	TodoListID uuid.UUID `json:"todolist_id" gorm:"type:uuid;default:uuid_generate_v4()"`
	TodoList   TodoList  `json:"-"`

	Files []TodoFile `json:"files"`

	Steps []TodoStep `json:"steps"`
}

type TodoFile struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	FileID   uuid.UUID `json:"file_id"`
	FileName string    `json:"file_name" gorm:"-"`
	File     File      `json:"-"`

	TodoID uuid.UUID `json:"todo_id"`
	Todo   Todo      `json:"-"`
}

type TodoStep struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	Name   string    `json:"name" gorm:"-"`
	Done   bool      `json:"done"`
	DoneAt time.Time `json:"done_at"`

	TodoID uuid.UUID `json:"todo_id"`
	Todo   Todo      `json:"-"`
}
