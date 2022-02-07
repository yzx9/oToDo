package dal

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoLists = make(map[uuid.UUID]entity.TodoList)

func InsertTodoList(todoList entity.TodoList) (entity.TodoList, error) {
	todoLists[todoList.ID] = todoList
	return todoList, nil
}

func GetTodoList(id uuid.UUID) (entity.TodoList, error) {
	for _, v := range todoLists {
		if v.ID == id {
			return v, nil
		}
	}

	return entity.TodoList{}, fmt.Errorf("todo list not found: %v", id)
}

func GetTodoLists(userId uuid.UUID) ([]entity.TodoList, error) {
	vec := make([]entity.TodoList, 0)
	for _, v := range todoLists {
		if v.UserID == userId {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func DeleteTodoList(todoListID uuid.UUID) (entity.TodoList, error) {
	todoList, ok := todoLists[todoListID]
	if !ok {
		return entity.TodoList{}, utils.NewErrorWithNotFound("todo list not found: %v", todoListID)
	}

	delete(todoLists, todoListID)
	return todoList, nil
}

func ExistTodoList(id uuid.UUID) bool {
	_, exist := todoLists[id]
	return exist
}
