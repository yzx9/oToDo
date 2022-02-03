package dal

import (
	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
)

// inject fake data
func init() {
	// User
	userID := uuid.MustParse("0c13da37-4593-4b2e-8163-1cbdb6e50830")
	users[userID] = entity.User{
		ID:       userID,
		Name:     "admin",
		Nickname: "Admin",
		Password: []byte{ // admin123
			0x92, 0x0E, 0xE3, 0xA9, 0xBE, 0xFC, 0x3E, 0xB3,
			0xB5, 0xB9, 0x79, 0x4B, 0xA9, 0xCE, 0x4D, 0xD3,
			0x04, 0x4B, 0x41, 0x39, 0x32, 0xD3, 0x4B, 0xDC,
			0xEB, 0x02, 0xDE, 0x90, 0x0A, 0xF2, 0x55, 0x36},
	}

	// Todo List
	todoListID := uuid.MustParse("5f5459d1-ffdb-40ce-9e05-02af49938a45")
	todoLists[todoListID] = entity.TodoList{
		ID:     todoListID,
		Name:   "To-Do",
		UserID: userID,
	}

	// Todo
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

	// File Template
	fileDestTemplateID := uuid.New()
	filePathTemplates[fileDestTemplateID] = entity.FilePathTemplate{
		ID:        fileDestTemplateID,
		Available: true,
		Type:      string(entity.FilePathTemplateTypeDest),
		Template:  "./file/:date/:filename.:ext",
	}

	fileServerTemplateID := uuid.New()
	filePathTemplates[fileServerTemplateID] = entity.FilePathTemplate{
		ID:        fileServerTemplateID,
		Available: true,
		Type:      string(entity.FilePathTemplateTypeServer),
		Template:  "//localhost:8080/file/:filename.:ext",
	}
}
