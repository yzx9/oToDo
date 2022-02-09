package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoList(userID uuid.UUID, todoListName string) (entity.TodoList, error) {
	list := entity.TodoList{
		ID:        uuid.New(),
		Name:      todoListName,
		Deletable: true,
		UserID:    userID,
	}
	if err := dal.InsertTodoList(list); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to create todo list: %v", todoListName)
	}

	return list, nil
}

func GetTodoList(userID, todoListID uuid.UUID) (entity.TodoList, error) {
	return OwnTodoList(userID, todoListID)
}

func GetTodoLists(userID uuid.UUID) ([]entity.TodoList, error) {
	vec, err := dal.GetTodoLists(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get user: %w", err)
	}

	// TODO: shared todo list

	return vec, nil
}

func DeleteTodoList(userID, todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	if !todoList.Deletable {
		return entity.TodoList{}, fmt.Errorf("todo list not deletable: %v", todoListID)
	}

	if err = dal.DeleteTodoList(todoListID); err != nil {
		return entity.TodoList{}, err
	}

	return todoList, nil
}

// Verify permission
func OwnTodoList(userID, todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, err := dal.GetTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %v", todoListID)
	}

	// TODO: shared todo list
	if todoList.UserID != userID {
		return entity.TodoList{}, utils.NewErrorWithForbidden("unable to handle non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}
