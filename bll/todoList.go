package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func CreateTodoList(userID uuid.UUID, todoListName string) (entity.TodoList, error) {
	return dal.InsertTodoList(entity.TodoList{
		ID:        uuid.New(),
		Name:      todoListName,
		Deletable: true,
		UserID:    userID,
	})
}

func GetTodoLists(userID uuid.UUID) ([]entity.TodoList, error) {
	vec, err := dal.GetTodoLists(userID)
	if err != nil {
		return nil, err
	}

	// TODO: shared todo list

	return vec, nil
}

func DeleteTodoList(todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, err := dal.GetTodoList(todoListID)
	if err != nil {
		return entity.TodoList{}, err
	}

	if !todoList.Deletable {
		return entity.TodoList{}, fmt.Errorf("todo list not deletable: %v", todoListID)
	}

	return dal.DeleteTodoList(todoListID)
}
