package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func GetTodoLists(userId string) ([]entity.TodoList, error) {
	id, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", id)
	}

	vec, err := dal.GetTodoLists(id)
	if err != nil {
		return nil, err
	}

	// TODO: add shared todo list

	return vec, nil
}
