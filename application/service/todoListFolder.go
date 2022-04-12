package service

import (
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todolist"
)

func CreateTodoListFolder(userID int64, folder dto.NewTodoListFolder) (dto.TodoListFolder, error) {
	entity := folder.ToEntity()
	if err := todolist.CreateTodoListFolder(userID, &entity); err != nil {
		return dto.TodoListFolder{}, err
	}

	return dto.TodoListFolder{}.FromEntity(entity), nil
}

func DeleteTodoListFolder(userID, todoListFolderID int64) (todolist.TodoListFolder, error) {
	return todolist.DeleteTodoListFolder(userID, todoListFolderID)
}
