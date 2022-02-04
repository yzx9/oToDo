package dal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todos = make(map[uuid.UUID]entity.Todo)

func InsertTodo(todo entity.Todo) (entity.Todo, error) {
	todos[todo.ID] = todo
	return todo, nil
}

func GetTodos(todoListID uuid.UUID) ([]entity.Todo, error) {
	vec := make([]entity.Todo, 0, len(todos))
	for _, v := range todos {
		if v.TodoListID == todoListID {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func GetTodo(id uuid.UUID) (entity.Todo, error) {
	todo, ok := todos[id]
	if !ok {
		return entity.Todo{}, utils.NewErrorWithHttpStatus(http.StatusNotFound, "todo not found: %v", id)
	}

	return todo, nil
}

func UpdateTodo(todo entity.Todo) (entity.Todo, error) {
	_, exists := todos[todo.ID]
	if !exists {
		return entity.Todo{}, utils.NewErrorWithHttpStatus(http.StatusNotFound, "todo not found: %v", todo.ID)
	}

	todos[todo.ID] = todo
	return todo, nil
}

func DeleteTodo(id uuid.UUID) (entity.Todo, error) {
	todo, ok := todos[id]
	if !ok {
		return entity.Todo{}, utils.NewErrorWithHttpStatus(http.StatusNotFound, "todo not found: %v", id)
	}

	delete(todos, id)
	return todo, nil
}
