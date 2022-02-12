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
	return filterTodos(func(t *entity.Todo) bool {
		return t.TodoListID == todoListID
	}), nil
}

func GetImportantTodos(userID string) ([]entity.Todo, error) {
	vec := filterTodos(func(t *entity.Todo) bool {
		return t.UserID == userID && t.Importance
	})

	return vec, nil
}

func GetPlanedTodos(userID string) ([]entity.Todo, error) {
	vec := filterTodos(func(t *entity.Todo) bool {
		return t.UserID == userID && !t.Deadline.IsZero()
	})

	sort.Slice(vec, func(i, j int) bool {
		return vec[i].Deadline.After(vec[j].Deadline)
	})

	return vec, nil
}

func GetNotNotifiedTodos(userID string) ([]entity.Todo, error) {
	return filterTodos(func(t *entity.Todo) bool {
		return t.UserID == userID && !t.Done && !t.Notified
	}), nil
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

func DeleteTodos(todoListID string) error {
	for _, v := range todos {
		if v.TodoListID == todoListID {
			delete(todos, v.ID)
		}
	}
	return nil
}

func filterTodos(filter func(*entity.Todo) bool) []entity.Todo {
	vec := make([]entity.Todo, 0)
	for _, v := range todos {
		if filter(&v) {
			files, _ := GetTodoFiles(v.ID)
			v.Files = files

			steps, _ := GetTodoSteps(v.ID)
			v.Steps = steps

			vec = append(vec, v)
		}
	}

	return vec
}
