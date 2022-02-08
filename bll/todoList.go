package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodoList(userID uuid.UUID, todoListName string) (entity.TodoList, error) {
	return dal.InsertTodoList(entity.TodoList{
		ID:        uuid.New(),
		Name:      todoListName,
		Deletable: true,
		UserID:    userID,
	})
}

func GetTodoList(userID uuid.UUID, todoListID uuid.UUID) (entity.TodoList, error) {
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

func DeleteTodoList(userID uuid.UUID, todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	if !todoList.Deletable {
		return entity.TodoList{}, fmt.Errorf("todo list not deletable: %v", todoListID)
	}

	return dal.DeleteTodoList(todoListID)
}

// Verify permission
func OwnTodoList(userID uuid.UUID, todoListID uuid.UUID) (entity.TodoList, error) {
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
