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
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FileName   string
	FilePath   string
	AccessType string    // FileAccessType
	RelatedID  uuid.UUID // depend on access type,
	CreatedAt  time.Time

	FileDestTemplateID uuid.UUID
	FileDestTemplate   FilePathTemplate `json:"-"`

	FileServerTemplateID uuid.UUID
	FileServerTemplate   FilePathTemplate `json:"-"`
}

type FilePathTemplateType string

const (
	FilePathTemplateTypeDest   FilePathTemplateType = "dest"
	FilePathTemplateTypeServer FilePathTemplateType = "server"
)

// Template
//
// Should replace :filename with real file name
type FilePathTemplate struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Available bool
	Type      string
	Template  string // FilePathTemplateType
	Host      string
	CreatedAt time.Time
}
