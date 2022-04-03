package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func CreateTodoListFolder(userID int64, folder *entity.TodoListFolder) error {
	folder.UserID = userID
	if err := repository.InsertTodoListFolder(folder); err != nil {
		return fmt.Errorf("fails to create todo list folder: %w", err)
	}

	return nil
}

func GetTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	return OwnTodoListFolder(userID, todoListFolderID)
}

func GetTodoListFolders(userID int64) ([]entity.TodoListFolder, error) {
	vec, err := repository.SelectTodoListFolders(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todo list folder: %w", err)
	}

	return vec, nil
}

func DeleteTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	write := func(err error) (entity.TodoListFolder, error) {
		return entity.TodoListFolder{}, err
	}

	folder, err := OwnTodoListFolder(userID, todoListFolderID)
	if err != nil {
		return write(err)
	}

	// TODO[feat] Whether to cascade delete todo lists
	// Cascade delete todo lists
	if _, err = repository.DeleteTodoListsByFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to cascade delete todo lists: %w", err))
	}

	if err = repository.DeleteTodoListFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to delete todo list folder: %w", err))
	}

	return folder, nil
}

// Verify permission
func OwnTodoListFolder(userID, todoListFolderID int64) (entity.TodoListFolder, error) {
	todoListFolder, err := repository.SelectTodoListFolder(todoListFolderID)
	if err != nil {
		return entity.TodoListFolder{}, fmt.Errorf("fails to get todo list folder: %v", todoListFolderID)
	}

	if todoListFolder.UserID != userID {
		return entity.TodoListFolder{}, util.NewErrorWithForbidden("unable to handle non-owned todo list folder: %v", todoListFolderID)
	}

	return todoListFolder, nil
}
