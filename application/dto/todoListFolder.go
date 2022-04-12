package dto

import "github.com/yzx9/otodo/domain/todolist"

type NewTodoListFolder struct {
	Name string `json:"name"`
}

func (f NewTodoListFolder) ToEntity() todolist.TodoListFolder {
	return todolist.TodoListFolder{
		Name: f.Name,
	}
}

type TodoListFolder struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (TodoListFolder) FromEntity(entity todolist.TodoListFolder) TodoListFolder {
	return TodoListFolder{
		ID:   entity.ID,
		Name: entity.Name,
	}
}
