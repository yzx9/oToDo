package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
)

func CreateTodo(userID int64, newTodo dto.NewTodo) (dto.Todo, error) {
	entity := newTodo.ToEntity()
	if err := todo.CreateTodo(userID, &entity); err != nil {
		return dto.Todo{}, err
	}

	return dto.Todo{}.FromEntity(entity), nil
}

func UpdateTodo(userID int64, t dto.Todo) (dto.Todo, error) {
	entity := t.ToEntity()
	if err := todo.UpdateTodo(userID, &entity); err != nil {
		return dto.Todo{}, nil
	}

	return dto.Todo{}.FromEntity(entity), nil
}

func DeleteTodo(userID int64, todoID int64) (todo.Todo, error) {
	return todo.DeleteTodo(userID, todoID)
}

func CreateTodoStep(userID int64, dto dto.NewTodoStep) (todo.TodoStep, error) {
	return todo.CreateTodoStep(userID, dto.TodoID, dto.Name)
}

func UpdateTodoStep(userID int64, step dto.TodoStep) (todo.TodoStep, error) {
	entity := step.ToEntity()
	if err := todo.UpdateTodoStep(userID, &entity); err != nil {
		return todo.TodoStep{}, nil
	}

	return entity, nil
}

func DeleteTodoStep(userID, todoID, todoStepID int64) (todo.TodoStep, error) {
	return todo.DeleteTodoStep(userID, todoID, todoStepID)
}

func GetTodo(userID, todoID int64) (todo.Todo, error) {
	entity, err := todo.OwnTodo(userID, todoID)
	if err != nil {
		return todo.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	return entity, nil
}

func GetTodos(userID, todoListID int64) ([]todo.Todo, error) {
	if _, err := todolist.OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	return ForceGetTodos(todoListID)
}

func ForceGetTodos(todoListID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllByTodoList(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetImportantTodos(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllImportantOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get important todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodos(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllPlanedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodos(userID int64) ([]todo.Todo, error) {
	todos, err := TodoRepository.FindAllNotNotifiedOnesByUser(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}
