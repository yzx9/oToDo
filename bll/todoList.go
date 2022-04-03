package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
)

func CreateTodoList(userID int64, todoList *entity.TodoList) error {
	todoList.IsBasic = false
	todoList.UserID = userID
	todoList.TodoListFolderID = 0
	if err := repository.InsertTodoList(todoList); err != nil {
		return fmt.Errorf("fails to create todo list: %w", err)
	}

	return nil
}

func GetTodoList(userID, todoListID int64) (entity.TodoList, error) {
	return OwnOrSharedTodoList(userID, todoListID)
}

func ForceGetTodoList(todoListID int64) (entity.TodoList, error) {
	list, err := repository.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %w", err)
	}

	return list, nil
}

func GetTodoLists(userID int64) ([]entity.TodoList, error) {
	vec, err := repository.SelectTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user todo lists: %w", err)
	}

	shared, err := repository.SelectSharedTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user shared todo lists: %w", err)
	}

	vec = append(vec, shared...)
	return vec, nil
}

func UpdateTodoList(userID int64, todoList *entity.TodoList) error {
	oldTodoList, err := OwnTodoList(userID, todoList.ID)
	if err != nil {
		return err
	}

	if oldTodoList.IsBasic {
		return util.NewErrorWithForbidden("unable to update basic todo list")
	}

	if err := repository.SaveTodoList(todoList); err != nil {
		return fmt.Errorf("fails to update todo list: %w", err)
	}

	return nil
}

func DeleteTodoList(userID, todoListID int64) (entity.TodoList, error) {
	// only allow delete by owner, not shared users
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	// disable to delete basic todo list
	if todoList.IsBasic {
		return entity.TodoList{}, util.NewErrorWithPreconditionFailed("unable to delete basic todo list: %v", todoListID)
	}

	// cascade delete todos
	if _, err = repository.DeleteTodos(todoListID); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to cascade delete todos: %w", err)
	}

	if err = repository.DeleteTodoList(todoListID); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to delete todo list: %w", err)
	}

	return todoList, nil
}

// owner
func OwnTodoList(userID, todoListID int64) (entity.TodoList, error) {
	todoList, err := repository.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		return entity.TodoList{}, util.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}
