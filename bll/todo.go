package bll

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodo(todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo, err := dal.InsertTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func GetTodo(id uuid.UUID) (entity.Todo, error) {
	todo, err := dal.GetTodo(id)
	if err != nil {
		return entity.Todo{}, err
	}

	// Fill files
	for _, file := range todo.Files {
		filePath, err := GetFilePath(file.ID)
		if err != nil {
			file.Valid = false
			continue
		}
		file.FilePath = filePath

		record, err := GetFile(file.ID)
		if err != nil {
			file.Valid = false
			continue
		}
		file.FileName = record.FileName
	}

	return todo, nil
}

func GetTodos(todoListID uuid.UUID) ([]entity.Todo, error) {
	if !dal.ExistTodoList(todoListID) {
		return nil, utils.NewErrorWithNotFound("todo list not found: %v", todoListID)
	}

	todos, err := dal.GetTodos(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func UpdateTodo(todo entity.Todo) (entity.Todo, error) {
	todo, err := dal.InsertTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func DeleteTodo(id uuid.UUID) (entity.Todo, error) {
	return dal.DeleteTodo(id)
}
