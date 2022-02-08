package entity

import (
	"time"

	"github.com/google/uuid"
)

type FileAccessType string

const (
	FileTypePublic FileAccessType = "public" // set RelatedID to empty
	FileTypeTodo   FileAccessType = "todo"   // set RelatedID to TodoID
)

type File struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	AccessType string    `json:"access_type"` // FileAccessType
	RelatedID  uuid.UUID `json:"related_id"`  // Depend on access type
	CreatedAt  time.Time `json:"created_at"`
}
