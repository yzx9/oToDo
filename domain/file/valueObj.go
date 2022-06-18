package file

type FileAccessType int8

const (
	FileTypePublic FileAccessType = 10*iota + 1 // set RelatedID to empty
	FileTypeTodo                                // set RelatedID to TodoID
)
