package dal

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todos = make(map[uuid.UUID]entity.Todo)

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

func AddTodo(todo entity.Todo) error {
	todos[todo.ID] = todo
	return nil
}