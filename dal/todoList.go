package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoLists = make(map[string]entity.TodoList)

func InsertTodoList(todoList entity.TodoList) error {
	todoLists[todoList.ID] = todoList
	return nil
}

func GetTodoList(id string) (entity.TodoList, error) {
	for _, v := range todoLists {
		if v.ID == id {
			return v, nil
		}
	}

	return entity.TodoList{}, utils.NewErrorWithNotFound("todo list not found: %v", id)
}

func GetTodoLists(userId string) ([]entity.TodoList, error) {
	vec := make([]entity.TodoList, 0)
	for _, v := range todoLists {
		if v.UserID == userId {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func DeleteTodoList(todoListID string) error {
	_, ok := todoLists[todoListID]
	if !ok {
		return utils.NewErrorWithNotFound("todo list not found: %v", todoListID)
	}

	delete(todoLists, todoListID)
	return nil
}

func ExistTodoList(id string) bool {
	_, exist := todoLists[id]
	return exist
}
