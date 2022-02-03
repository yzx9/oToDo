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
	RelatedID  uuid.UUID `json:"related_id"`  // depend on access type,
	CreatedAt  time.Time `json:"created_at"`

	FileDestTemplateID uuid.UUID        `json:"file_dest_template_id"`
	FileDestTemplate   FilePathTemplate `json:"-"`

	FileServerTemplateID uuid.UUID        `json:"file_server_template_id"`
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
	Available bool      `json:"available"`
	Type      string    `json:"type"`
	Template  string    `json:"template"` // FilePathTemplateType
	Host      string    `json:"host"`
	CreatedAt time.Time `json:"created_at"`
}
