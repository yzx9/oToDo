package dal

import (
	"github.com/google/uuid"

	"github.com/yzx9/otodo/entity"
)

func init() {
	// inject fake data

	userID := uuid.MustParse("0c13da37-4593-4b2e-8163-1cbdb6e50830")
	users[userID] = entity.User{
		ID:   userID,
		Name: "Admin",
	}

	todoListID := uuid.MustParse("5f5459d1-ffdb-40ce-9e05-02af49938a45")
	todoLists[todoListID] = entity.TodoList{
		ID:     todoListID,
		Name:   "Super",
		UserID: userID,
	}

	AddTodo(entity.Todo{
		ID:         uuid.MustParse("32acb375-e9dc-473e-8f5f-8826f7783c1d"),
		Title:      "Hello, World!",
		UserID:     userID,
		TodoListID: todoListID,
	})

	AddTodo(entity.Todo{
		ID:         uuid.MustParse("343dc2ce-1fbc-43ad-98d6-9cac1c67f2a6"),
		Title:      "你好，世界！",
		UserID:     userID,
		TodoListID: todoListID,
	})
}
