package entity

type FileAccessType int

const (
	FileTypePublic FileAccessType = 10*iota + 1 // set RelatedID to empty
	FileTypeTodo                                // set RelatedID to TodoID
)

type File struct {
	Entity

	FileName     string `json:"fileName"`
	FileServerID string `json:"-" gorm:"size:15"`
	FilePath     string `json:"-" gorm:"size:128"`
	AccessType   int8   `json:"accessType"` // FileAccessType
	RelatedID    int64  `json:"relatedID"`  // Depend on access type
}
