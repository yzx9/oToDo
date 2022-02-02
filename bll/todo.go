package bll

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func GetTodos(todoListID string) ([]entity.Todo, error) {
	id, err := uuid.Parse(todoListID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", todoListID)
	}

	return dal.GetTodos(id), nil
}

func GetTodo(id string) (entity.Todo, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return entity.Todo{}, errors.New("invalid uuid")
	}

	todo, err := dal.GetTodo(uuid)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}
