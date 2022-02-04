package bll

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func GetTodos(todoListID string) ([]entity.Todo, error) {
	id, err := uuid.Parse(todoListID)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", todoListID)
	}

	if !dal.ExistTodoList(id) {
		return nil, utils.NewErrorWithHttpStatus(http.StatusNotFound, "todo list not found: %v", todoListID)
	}

	todos, err := dal.GetTodos(id)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
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

func AddTodo(todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo, err := dal.AddTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func UpdateTodo(todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo, err := dal.AddTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}
