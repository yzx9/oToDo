package dal

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
)

var todos = make(map[uuid.UUID]entity.Todo)

func GetTodos(todoListID uuid.UUID) []entity.Todo {
	vec := make([]entity.Todo, 0, len(todos))
	for _, v := range todos {
		if v.TodoListID == todoListID {
			vec = append(vec, v)
		}
	}
	return vec
}

func GetTodo(id uuid.UUID) (entity.Todo, error) {
	todo, ok := todos[id]
	if !ok {
		return entity.Todo{}, errors.New("id not exist")
	}
	return todo, nil
}

func AddTodo(todo entity.Todo) error {
	todos[todo.ID] = todo
	return nil
}
