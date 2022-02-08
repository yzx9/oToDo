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
	todoList, err := dal.GetTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to get todo list: %w", err)
	}

	if todoList.UserID != userID {
		return entity.TodoList{}, utils.NewErrorWithForbidden("unable to get non-owned todo list: %v", todoListID)
	}

	return todoList, nil
}

func GetTodoLists(userID uuid.UUID) ([]entity.TodoList, error) {
	vec, err := dal.GetTodoLists(userID)
	if err != nil {
		return nil, err
	}

	// TODO: shared todo list

	return vec, nil
}

func DeleteTodoList(userID uuid.UUID, todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, err := dal.GetTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	if todoList.UserID != userID {
		return entity.TodoList{}, utils.NewErrorWithForbidden("unable to delete non-owned todo list: %v", todoListID)
	}

	if !todoList.Deletable {
		return entity.TodoList{}, fmt.Errorf("todo list not deletable: %v", todoListID)
	}

	return dal.DeleteTodoList(todoListID)
}
