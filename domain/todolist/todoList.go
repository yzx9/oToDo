package todolist

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

type TodoList struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	IsBasic   bool
	IsSharing bool

	UserID int64

	TodoListFolderID int64

	SharedUsers []*int64 // TODO
}

func CreateTodoList(userID int64, todoList *TodoList) error {
	todoList.ID = 0
	todoList.IsBasic = false
	todoList.UserID = userID
	todoList.TodoListFolderID = 0
	if err := TodoListRepository.Save(todoList); err != nil {
		return fmt.Errorf("fails to create todo list: %w", err)
	}

	return nil
}

func GetTodoList(userID, todoListID int64) (TodoList, error) {
	return OwnOrSharedTodoList(userID, todoListID)
}

func UpdateTodoList(userID int64, todoList *TodoList) error {
	oldTodoList, err := OwnTodoList(userID, todoList.ID)
	if err != nil {
		return err
	}

	if oldTodoList.IsBasic {
		return util.NewErrorWithForbidden("unable to update basic todo list")
	}

	if err := TodoListRepository.Save(todoList); err != nil {
		return fmt.Errorf("fails to update todo list: %w", err)
	}

	return nil
}

func DeleteTodoList(userID, todoListID int64) (TodoList, error) {
	// only allow delete by owner, not shared users
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return TodoList{}, err
	}

	// disable to delete basic todo list
	if todoList.IsBasic {
		return TodoList{}, util.NewErrorWithPreconditionFailed("unable to delete basic todo list: %v", todoListID)
	}

	// cascade delete todos
	if _, err = TodoRepository.DeleteAllByTodoList(todoListID); err != nil {
		return TodoList{}, fmt.Errorf("fails to cascade delete todos: %w", err)
	}

	if err = TodoListRepository.Delete(todoListID); err != nil {
		return TodoList{}, fmt.Errorf("fails to delete todo list: %w", err)
	}

	return todoList, nil
}

// owner
func OwnTodoList(userID, todoListID int64) (TodoList, error) {
	todoList, err := TodoListRepository.Find(todoListID)
	if err != nil {
		return TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		return TodoList{}, util.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}
