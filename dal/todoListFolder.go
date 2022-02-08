package dal

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoListFolders = make(map[uuid.UUID]entity.TodoListFolder)

func InsertTodoListFolder(todoListFolder entity.TodoListFolder) (entity.TodoListFolder, error) {
	todoListFolders[todoListFolder.ID] = todoListFolder
	return todoListFolder, nil
}

func GetTodoListFolder(todoListFolderID uuid.UUID) (entity.TodoListFolder, error) {
	for _, v := range todoListFolders {
		if v.ID == todoListFolderID {
			return v, nil
		}
	}

	return entity.TodoListFolder{}, fmt.Errorf("todo list folder not found: %v", todoListFolderID)
}

func GetTodoListFolders(userId uuid.UUID) ([]entity.TodoListFolder, error) {
	vec := make([]entity.TodoListFolder, 0)
	for _, v := range todoListFolders {
		if v.UserID == userId {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func DeleteTodoListFolder(todoListFolderID uuid.UUID) (entity.TodoListFolder, error) {
	todoListFolder, ok := todoListFolders[todoListFolderID]
	if !ok {
		return entity.TodoListFolder{}, utils.NewErrorWithNotFound("todo list folder not found: %v", todoListFolderID)
	}

	delete(todoListFolders, todoListFolderID)
	return todoListFolder, nil
}

func ExistTodoListFolder(id uuid.UUID) bool {
	_, exist := todoLists[id]
	return exist
}
