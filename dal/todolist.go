package dal

import (
	"github.com/google/uuid"

	"github.com/yzx9/otodo/entity"
)

var todoLists = make(map[uuid.UUID]entity.TodoList)

func GetTodoLists(userId uuid.UUID) ([]entity.TodoList, error) {
	vec := make([]entity.TodoList, 0, len(todoLists))
	for _, v := range todoLists {
		if v.UserID == userId {
			vec = append(vec, v)
		}
	}

	return vec, nil
}
