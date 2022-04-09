package service

import (
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/domain/user"
)

var UserRepository userRepository
var TodoRepository todoRepository
var TodoListRepository todoListRepository
var TodoListFolderRepository todoListFolderRepository
var TodoListSharingRepository todoListSharingRepository
var SharingRepository sharingRepository

type userRepository interface {
	Find(id int64) (user.User, error)
}

type todoRepository interface {
	FindAllByTodoList(todoListID int64) ([]todo.Todo, error)

	FindAllImportantOnesByUser(userID int64) ([]todo.Todo, error)

	FindAllPlanedOnesByUser(userID int64) ([]todo.Todo, error)

	FindAllNotNotifiedOnesByUser(userID int64) ([]todo.Todo, error)
}

type todoListRepository interface {
	Find(id int64) (todolist.TodoList, error)

	FindByUserWithMenuFormat(userID int64) ([]dto.TodoListMenuItem, error)

	FindAllByUser(userID int64) ([]todolist.TodoList, error)

	FindAllSharedByUser(userID int64) ([]todolist.TodoList, error)
}

type todoListFolderRepository interface {
	FindAllByUser(userId int64) ([]todolist.TodoListFolder, error)
}

type sharingRepository interface {
	FindAllActive(userID int64, sharingType todolist.SharingType) ([]todolist.Sharing, error)
}

type todoListSharingRepository interface {
	FindAllSharedUsers(todoListID int64) ([]user.User, error)
}
