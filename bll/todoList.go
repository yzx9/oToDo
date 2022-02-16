package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func CreateTodoList(userID int64, todoList *entity.TodoList) error {
	todoList.Deletable = false
	todoList.UserID = userID
	todoList.TodoListFolderID = 0
	if err := dal.InsertTodoList(todoList); err != nil {
		return fmt.Errorf("fails to create todo list: %w", err)
	}

	return nil
}

func SelectTodoList(userID, todoListID int64) (entity.TodoList, error) {
	return OwnOrSharedTodoList(userID, todoListID)
}

func SelectTodoLists(userID int64) ([]entity.TodoList, error) {
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

func DeleteTodoList(userID, todoListID int64) (entity.TodoList, error) {
	// only allow delete by owner, not shared users
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	// check if deletable
	if !todoList.Deletable {
		return entity.TodoList{}, util.NewErrorWithPreconditionFailed("todo list not deletable: %v", todoListID)
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

// owner
func OwnTodoList(userID, todoListID int64) (entity.TodoList, error) {
	todoList, err := dal.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	if todoList.UserID != userID {
		return entity.TodoList{}, util.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}
