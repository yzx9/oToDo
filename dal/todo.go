package dal

import (
	"sort"

	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todos = make(map[string]entity.Todo)

func InsertTodo(todo entity.Todo) error {
	todos[todo.ID] = todo
	return nil
}

func GetTodo(id string) (entity.Todo, error) {
	todo, ok := todos[id]
	if !ok {
		return entity.Todo{}, utils.NewErrorWithNotFound("todo not found: %v", id)
	}

	return todo, nil
}

func GetTodos(todoListID string) ([]entity.Todo, error) {
	vec := make([]entity.Todo, 0, len(todos))
	for _, v := range todos {
		if v.TodoListID == todoListID {
			files, _ := GetTodoFiles(v.ID)
			v.Files = files

			steps, _ := GetTodoSteps(v.ID)
			v.Steps = steps

			vec = append(vec, v)
		}
	}

	return vec, nil
}

func GetPlanedTodos(userID string) ([]entity.Todo, error) {
	vec := make([]entity.Todo, 0, len(todos))
	for _, v := range todos {
		if v.UserID == userID && !v.Deadline.IsZero() {
			files, _ := GetTodoFiles(v.ID)
			v.Files = files

			steps, _ := GetTodoSteps(v.ID)
			v.Steps = steps

			vec = append(vec, v)
		}
	}

	sort.Slice(vec, func(i, j int) bool {
		return vec[i].Deadline.After(vec[j].Deadline)
	})

	return vec, nil
}

func UpdateTodo(todo entity.Todo) error {
	_, exists := todos[todo.ID]
	if !exists {
		return utils.NewErrorWithNotFound("todo not found: %v", todo.ID)
	}

	todos[todo.ID] = todo
	return nil
}

func DeleteTodo(id string) error {
	_, ok := todos[id]
	if !ok {
		return utils.NewErrorWithNotFound("todo not found: %v", id)
	}

	delete(todos, id)
	return nil
}
