package service

import (
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/domain/sharing"
	"github.com/yzx9/otodo/domain/todo"
)

var UserRepository interface {
	Find(id int64) (identity.User, error)
}

var TodoRepository interface {
	Find(id int64) (todo.Todo, error)
	FindAllByTodoList(todoListID int64) ([]todo.Todo, error)
	FindAllInBasicTodoList(userID int64) ([]todo.Todo, error)
	FindAllImportantOnesByUser(userID int64) ([]todo.Todo, error)
	FindAllPlanedOnesByUser(userID int64) ([]todo.Todo, error)
	FindAllNotNotifiedOnesByUser(userID int64) ([]todo.Todo, error)
}

var TodoListRepository interface {
	Find(id int64) (todo.TodoList, error)
	FindAllByUser(userID int64) ([]todo.TodoList, error)
	FindAllSharedByUser(userID int64) ([]todo.TodoList, error)
}

var TodoListFolderRepository interface {
	FindAllByUser(userId int64) ([]todo.TodoListFolder, error)
}

var SharingRepository interface {
	FindAllActive(userID int64, sharingType sharing.SharingType) ([]sharing.Sharing, error)
}

var TodoListSharingRepository interface {
	FindAllSharedUsers(todoListID int64) ([]identity.User, error)
}
