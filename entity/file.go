package entity

type FileAccessType string

const (
	FileTypePublic FileAccessType = "public" // set RelatedID to empty
	FileTypeTodo   FileAccessType = "todo"   // set RelatedID to TodoID
)

type File struct {
	Entity

	FileName     string `json:"file_name"`
	FileServerID string `json:"-" gorm:"size:15"`
	FilePath     string `json:"file_path"`
	AccessType   string `json:"access_type" gorm:"size:8"`        // FileAccessType
	RelatedID    string `json:"related_id" gorm:"type:char(36);"` // Depend on access type
}
