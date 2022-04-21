package service

import (
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todo"
)

func CreateTodoListFolder(userID int64, folder dto.NewTodoListFolder) (dto.TodoListFolder, error) {
	entity := folder.ToEntity()
	if err := todo.CreateTodoListFolder(userID, &entity); err != nil {
		return dto.TodoListFolder{}, err
	}

	return dto.TodoListFolder{}.FromEntity(entity), nil
}

func DeleteTodoListFolder(userID, todoListFolderID int64) (todo.TodoListFolder, error) {
	return todo.DeleteTodoListFolder(userID, todoListFolderID)
}
