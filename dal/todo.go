package dal

import (
	"errors"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
)

var todos = make(map[string]entity.Todo)

func GetTodos() []entity.Todo {
	vec := make([]entity.Todo, 0, len(todos))
	for _, v := range todos {
		vec = append(vec, v)
	}
	return vec
}

func GetTodo(id string) (entity.Todo, error) {
	todo, ok := todos[id]
	if !ok {
		return entity.Todo{}, errors.New("id not exist")
	}
	return todo, nil
}

func AddTodo(todo entity.Todo) error {
	todos[todo.ID.String()] = todo
	return nil
}

func init() {
	var fake []entity.Todo

	fake = append(fake, entity.Todo{
		ID:    uuid.MustParse("32acb375-e9dc-473e-8f5f-8826f7783c1d"),
		Title: "Hello, World!",
	})

	fake = append(fake, entity.Todo{
		ID:    uuid.MustParse("343dc2ce-1fbc-43ad-98d6-9cac1c67f2a6"),
		Title: "你好，世界！",
	})

	for _, todo := range fake {
		todos[todo.ID.String()] = todo
	}
}
