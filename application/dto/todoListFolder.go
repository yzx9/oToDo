package dto

import "github.com/yzx9/otodo/domain/todo"

type NewTodoListFolder struct {
	Name string `json:"name"`
}

func (f NewTodoListFolder) ToEntity() todo.TodoListFolder {
	return todo.TodoListFolder{
		Name: f.Name,
	}
}

type TodoListFolder struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (TodoListFolder) FromEntity(entity todo.TodoListFolder) TodoListFolder {
	return TodoListFolder{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
