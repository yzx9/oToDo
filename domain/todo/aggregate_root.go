package todo

// TODO Create Aggregate Root
var TodoRepo TodoRepository

type TodoRepository interface {
	SaveFile(todoID, fileID int64) error
}
