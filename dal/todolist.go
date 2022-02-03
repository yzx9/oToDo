package dal

import (
	"github.com/google/uuid"

	"github.com/yzx9/otodo/entity"
)

var todoLists = make(map[uuid.UUID]entity.TodoList)

func ExistTodoList(id uuid.UUID) bool {
	_, exist := todoLists[id]
	return exist
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
