package todolist

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

type TodoListFolder struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string

	UserID int64

	TodoLists []int64 // TODO
}

func CreateTodoListFolder(userID int64, folder *TodoListFolder) error {
	folder.UserID = userID
	if err := TodoListFolderRepository.Save(folder); err != nil {
		return fmt.Errorf("fails to create todo list folder: %w", err)
	}

	return nil
}

func DeleteTodoListFolder(userID, todoListFolderID int64) (TodoListFolder, error) {
	write := func(err error) (TodoListFolder, error) {
		return TodoListFolder{}, err
	}

	folder, err := OwnTodoListFolder(userID, todoListFolderID)
	if err != nil {
		return write(err)
	}

	// TODO[feat] Whether to cascade delete todo lists
	// Cascade delete todo lists
	if _, err = TodoListRepository.DeleteAllByFolder(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to cascade delete todo lists: %w", err))
	}

	if err = TodoListFolderRepository.Delete(todoListFolderID); err != nil {
		return write(fmt.Errorf("fails to delete todo list folder: %w", err))
	}

	return folder, nil
}

// Verify permission
func OwnTodoListFolder(userID, todoListFolderID int64) (TodoListFolder, error) {
	todoListFolder, err := TodoListFolderRepository.Find(todoListFolderID)
	if err != nil {
		return TodoListFolder{}, fmt.Errorf("fails to get todo list folder: %v", todoListFolderID)
	}

	if todoListFolder.UserID != userID {
		return TodoListFolder{}, util.NewErrorWithForbidden("unable to handle non-owned todo list folder: %v", todoListFolderID)
	}

	return todoListFolder, nil
}
