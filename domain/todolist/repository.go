package todolist

import "github.com/yzx9/otodo/infrastructure/repository"

var TodoRepository todoRepository
var TodoListRepository todoListRepository
var TodoListFolderRepository todoListFolderRepository
var SharingRepository sharingRepository
var TodoListSharingRepository todoListSharingRepository

type todoRepository interface {
	DeleteAllByTodoList(todoListID int64) (int64, error)
}

type todoListRepository interface {
	Save(todoList *repository.TodoList) error

	Delete(id int64) error

	DeleteAllByFolder(todoListFolderID int64) (int64, error)

	Find(id int64) (repository.TodoList, error)
}

type todoListFolderRepository interface {
	Save(todoListFolder *repository.TodoListFolder) error

	Delete(id int64) error

	Find(id int64) (repository.TodoListFolder, error)
}

type sharingRepository interface {
	Save(sharing *repository.Sharing) error

	DeleteAllByUserAndType(userID int64, sharingType repository.SharingType) (int64, error)

	Find(token string) (repository.Sharing, error)
}

type todoListSharingRepository interface {
	SaveSharedUser(userID, todoListID int64) error

	DeleteSharedUser(userID, todoListID int64) error

	ExistSharing(userID, todoListID int64) (bool, error)
}
