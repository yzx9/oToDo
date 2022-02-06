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

func CreateTodo(todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo, err := dal.InsertTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
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

	// Fill files
	for _, file := range todo.Files {
		filePath, err := GetFilePath(file.ID.String())
		if err != nil {
			file.Valid = false
			continue
		}
		file.FilePath = filePath

		record, err := GetFile(file.ID.String())
		if err != nil {
			file.Valid = false
			continue
		}
		file.FileName = record.FileName
	}

	return todo, nil
}

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

func UpdateTodo(todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.New()
	todo, err := dal.InsertTodo(todo)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func DeleteTodo(id string) (entity.Todo, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return entity.Todo{}, errors.New("invalid uuid")
	}

	todo, err := dal.DeleteTodo(uuid)
	if err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}
