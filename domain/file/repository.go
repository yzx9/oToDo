package file

var FileRepository fileRepository

type fileRepository interface {
	Save(file *File) error

	Find(id int64) (*File, error)
}
