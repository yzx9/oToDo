package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoList(userID string, todoListName string) (entity.TodoList, error) {
	list := entity.TodoList{
		Name:      todoListName,
		Deletable: true,
		UserID:    userID,
	}
	if err := dal.InsertTodoList(&list); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to create todo list: %v", todoListName)
	}

	return list, nil
}

func GetTodoList(userID, todoListID string) (entity.TodoList, error) {
	return OwnTodoList(userID, todoListID)
}

func GetTodoLists(userID string) ([]entity.TodoList, error) {
	vec, err := dal.SelectTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user: %w", err)
	}

	// TODO[feat] shared todo list

	return vec, nil
}

func DeleteTodoList(userID, todoListID string) (entity.TodoList, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	// TODO[feat] only allow delete by owner, not shared users
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

// Verify permission
func OwnTodoList(userID, todoListID string) (entity.TodoList, error) {
	todoList, err := dal.SelectTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	// TODO: shared todo list
	if todoList.UserID != userID {
		return entity.TodoList{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}
