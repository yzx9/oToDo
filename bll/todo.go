package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodo(userID uuid.UUID, todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo.UserID = userID // override user
	todo, err := dal.InsertTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func GetTodo(userID uuid.UUID, todoID uuid.UUID) (entity.Todo, error) {
	return OwnTodo(userID, todoID)
}

func GetTodos(userID uuid.UUID, todoListID uuid.UUID) ([]entity.Todo, error) {
	if _, err := OwnTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	todos, err := dal.GetTodos(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func UpdateTodo(userID uuid.UUID, todo entity.Todo) (entity.Todo, error) {
	oldTodo, err := OwnTodo(userID, todo.ID)
	if err != nil {
		return entity.Todo{}, err
	}

	if oldTodo.UserID != todo.UserID {
		return entity.Todo{}, fmt.Errorf("unable to update todo owner")
	}

	todo, err = dal.UpdateTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func DeleteTodo(userID uuid.UUID, todoID uuid.UUID) (entity.Todo, error) {
	if _, err := OwnTodo(userID, todoID); err != nil {
		return entity.Todo{}, err
	}

	return dal.DeleteTodo(todoID)
}

func OwnTodo(userID uuid.UUID, todoID uuid.UUID) (entity.Todo, error) {
	todo, err := dal.GetTodo(todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %v", todoID)
	}

	// TODO todo in shared todo list
	if todo.UserID != userID {
		return entity.Todo{}, utils.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID)
	}

	return todo, nil
}
