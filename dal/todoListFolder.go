package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoListFolders = make(map[string]entity.TodoListFolder)

func InsertTodoListFolder(todoListFolder entity.TodoListFolder) error {
	todoListFolders[todoListFolder.ID] = todoListFolder
	return nil
}

func GetTodoListFolder(todoListFolderID string) (entity.TodoListFolder, error) {
	for _, v := range todoListFolders {
		if v.ID == todoListFolderID {
			return v, nil
		}
	}

	return entity.TodoListFolder{}, utils.NewErrorWithNotFound("todo list folder not found: %v", todoListFolderID)
}

func GetTodoListFolders(userId string) ([]entity.TodoListFolder, error) {
	vec := make([]entity.TodoListFolder, 0)
	for _, v := range todoListFolders {
		if v.UserID == userId {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func DeleteTodoListFolder(todoListFolderID string) error {
	_, ok := todoListFolders[todoListFolderID]
	if !ok {
		return utils.NewErrorWithNotFound("todo list folder not found: %v", todoListFolderID)
	}

	delete(todoListFolders, todoListFolderID)

	// TODO[feat] Whether to cascade delete todo lists
	// Remove from todo list
	for i, v := range todoLists {
		if v.TodoListFolderID == todoListFolderID {
			v.TodoListFolderID = ""
			todoLists[i] = v
		}
	}

	return nil
}

func ExistTodoListFolder(id string) bool {
	_, exist := todoLists[id]
	return exist
}
