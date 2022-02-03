package entity

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	FileName  string
	FilePath  string
	CreatedAt time.Time

	FileDestTemplateID uuid.UUID
	FileDestTemplate   FilePathTemplate

	FileServerTemplateID uuid.UUID
	FileServerTemplate   FilePathTemplate
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
	Template  string
	Host      string
	CreatedAt time.Time
}
