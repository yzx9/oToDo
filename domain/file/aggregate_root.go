package file

// TODO Create Aggregate Root
var FileRepo FileRepository

var TodoRepo TodoRepository // TODO remove

type FileRepository interface {
	Save(file *File) error

	Find(id int64) (*File, error)
}

type TodoRepository interface {
	SaveFile(todoID, fileID int64) error
}
