package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoList(userID string, todoList *entity.TodoList) error {
	todoList.Deletable = false
	todoList.UserID = userID
	todoList.TodoListFolderID = ""
	if err := dal.InsertTodoList(todoList); err != nil {
		return fmt.Errorf("fails to create todo list: %w", err)
	}

	return nil
}

func SelectTodoList(userID, todoListID string) (entity.TodoList, error) {
	return OwnTodoList(userID, todoListID)
}

func SelectTodoLists(userID string) ([]entity.TodoList, error) {
	vec, err := dal.SelectTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user todo lists: %w", err)
	}

	shared, err := dal.SelectSharedTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user shared todo lists: %w", err)
	}

	vec = append(vec, shared...)
	return vec, nil
}

func DeleteTodoList(userID, todoListID string) (entity.TodoList, error) {
	// only allow delete by owner, not shared users
	todoList, err := dal.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %w", err)
	}

	if todoList.UserID != userID {
		return entity.TodoList{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	// check if deletable
	if !todoList.Deletable {
		return entity.TodoList{}, utils.NewErrorWithPreconditionFailed("todo list not deletable: %v", todoListID)
	}

	// cascade delete todos
	if _, err = dal.DeleteTodos(todoListID); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to cascade delete todos: %w", err)
	}

	if err = dal.DeleteTodoList(todoListID); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to delete todo list: %w", err)
	}

	return todoList, nil
}

// Verify permission, owner or shared user
func OwnTodoList(userID, todoListID string) (entity.TodoList, error) {
	todoList, err := dal.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		sharing, err := HasSharing(userID, todoListID)
		if err != nil {
			return entity.TodoList{}, fmt.Errorf("fails to get todo list sharing: %v", todoListID)
		}

		if !sharing {
			return entity.TodoList{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
		}
	}

	return todoList, nil
}
