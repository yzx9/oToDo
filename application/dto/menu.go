package dto

import "github.com/yzx9/otodo/infrastructure/repository"

type TodoListMenuItem struct {
	repository.TodoListMenuItem

	IsLeaf   bool               `json:"isLeaf"`
	Children []TodoListMenuItem `json:"children"`
}
