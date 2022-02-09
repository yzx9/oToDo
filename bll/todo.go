package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodo(userID uuid.UUID, todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo.UserID = userID // override user
	if err := dal.InsertTodo(todo); err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func GetTodo(userID, todoID uuid.UUID) (entity.Todo, error) {
	return OwnTodo(userID, todoID)
}

func GetTodos(userID, todoListID uuid.UUID) ([]entity.Todo, error) {
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

	if !oldTodo.Done && todo.Done {
		todo.DoneAt = time.Now()
	}

	if err = dal.UpdateTodo(todo); err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func DeleteTodo(userID, todoID uuid.UUID) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, err
	}

	err = dal.DeleteTodo(todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to delete todo: %v", todoID)
	}

	return todo, nil
}

func OwnTodo(userID, todoID uuid.UUID) (entity.Todo, error) {
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
