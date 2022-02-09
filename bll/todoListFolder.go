package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoListFolder(userID uuid.UUID, todoListFolderName string) (entity.TodoListFolder, error) {
	folder := entity.TodoListFolder{
		ID:     uuid.New(),
		Name:   todoListFolderName,
		UserID: userID,
	}
	if err := dal.InsertTodoListFolder(folder); err != nil {
		return entity.TodoListFolder{}, err
	}

	return folder, nil
}

func GetTodoListFolder(userID, todoListFolderID uuid.UUID) (entity.TodoListFolder, error) {
	return OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID uuid.UUID) ([]entity.TodoListFolder, error) {
	vec, err := dal.GetTodoListFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user: %w", err)
	}

	return vec, nil
}

func DeleteTodoListFolder(userID, todoListFolderID uuid.UUID) (entity.TodoListFolder, error) {
	// TODO this is a BUG
	todoListFolder, err := OwnTodoList(userID, todoListFolderID)
	if err != nil {
		return entity.TodoListFolder{}, err
	}

	if !todoListFolder.Deletable {
		return entity.TodoListFolder{}, fmt.Errorf("todo list folder not deletable: %v", todoListFolderID)
	}

	return entity.TodoListFolder{}, dal.DeleteTodoListFolder(todoListFolderID)
}

// Verify permission
func OwnTodoListFolder(userID, todoListFolderID uuid.UUID) (entity.TodoListFolder, error) {
	todoListFolder, err := dal.GetTodoListFolder(todoListFolderID)
	if err != nil {
		return entity.TodoListFolder{}, fmt.Errorf("fails to get todo list folder: %v", todoListFolderID)
	}

	if todoListFolder.UserID != userID {
		return entity.TodoListFolder{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list folder: %v", todoListFolderID)
	}

	return todoListFolder, nil
}
