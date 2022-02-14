package entity

type FileAccessType string

const (
	FileTypePublic FileAccessType = "public" // set RelatedID to empty
	FileTypeTodo   FileAccessType = "todo"   // set RelatedID to TodoID
)

type File struct {
	Entity

	FileName   string `json:"file_name"`
	FilePath   string `json:"file_path"`
	AccessType string `json:"access_type"` // FileAccessType
	RelatedID  string `json:"related_id"`  // Depend on access type
}
