package bll

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func GetTodos(todoListID string) ([]entity.Todo, error) {
	uuid, err := uuid.Parse(todoListID)
	if err != nil {
		return nil, errors.New("invalid uuid")
	}

	return dal.GetTodos(uuid), nil
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
