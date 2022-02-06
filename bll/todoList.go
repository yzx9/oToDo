package bll

import (
	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func GetTodoLists(userID uuid.UUID) ([]entity.TodoList, error) {
	vec, err := dal.GetTodoLists(userID)
	if err != nil {
		return nil, err
	}

	// TODO: add shared todo list

	return vec, nil
}
