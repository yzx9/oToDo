package file

// TODO Create Aggregate Root
var FileRepository fileRepository

var TodoFileRepository todoFileRepository // TODO move to todo domain

type fileRepository interface {
	Save(file *File) error

	Find(id int64) (*File, error)
}

type todoFileRepository interface {
	Save(todoID, fileID int64) error
}
