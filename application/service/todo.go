package service

import (
	"fmt"

	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/repository"
)

func GetTodo(userID, todoID int64) (todo.Todo, error) {
	entity, err := todo.OwnTodo(userID, todoID)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	return entity, nil
}

func GetTodos(userID, todoListID int64) ([]repository.Todo, error) {
	if _, err := todolist.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	return ForceGetTodos(todoListID)
}

func ForceGetTodos(todoListID int64) ([]repository.Todo, error) {
	todos, err := repository.TodoRepo.FindAllByTodoList(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetImportantTodos(userID int64) ([]todo.Todo, error) {
	todos, err := repository.TodoRepo.FindAllImportantOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get important todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodos(userID int64) ([]todo.Todo, error) {
	todos, err := repository.TodoRepo.FindAllPlanedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodos(userID int64) ([]todo.Todo, error) {
	todos, err := repository.TodoRepo.FindAllNotNotifiedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}
