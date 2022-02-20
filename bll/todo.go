package bll

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func CreateTodo(userID int64, todo *entity.Todo) error {
	_, err := OwnOrSharedTodoList(userID, todo.TodoListID)
	if err != nil {
		return fmt.Errorf("fails to get todo list: %w", err)
	}

	todo.UserID = userID // override user

	plan, err := CreateTodoRepeatPlan(todo.TodoRepeatPlan)
	if err != nil {
		return fmt.Errorf("fails to create todo repeat plan: %w", err)
	}
	todo.TodoRepeatPlanID = plan.ID

	if err := dal.InsertTodo(todo); err != nil {
		return fmt.Errorf("fails to create todo: %w", err)
	}

	return nil
}

func GetTodo(userID, todoID int64) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	return todo, nil
}

func GetTodos(userID, todoListID int64) ([]entity.Todo, error) {
	if _, err := OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	return ForceGetTodos(todoListID)
}

func ForceGetTodos(todoListID int64) ([]entity.Todo, error) {
	todos, err := dal.SelectTodos(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetImportantTodos(userID int64) ([]entity.Todo, error) {
	todos, err := dal.SelectImportantTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get important todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodos(userID int64) ([]entity.Todo, error) {
	todos, err := dal.SelectPlanedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodos(userID int64) ([]entity.Todo, error) {
	todos, err := dal.SelectNotNotifiedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}

func UpdateTodo(userID int64, todo *entity.Todo) error {
	// Limits
	oldTodo, err := OwnTodo(userID, todo.ID)
	if err != nil {
		return err
	}

	todo.CreatedAt = oldTodo.CreatedAt
	todo.UserID = oldTodo.UserID
	todo.Files = oldTodo.Files
	todo.Steps = oldTodo.Steps
	todo.NextID = oldTodo.NextID

	if !oldTodo.Done && todo.Done {
		t := time.Now()
		todo.DoneAt = &t

		// Create Repeat Todo If Need
		if todo.NextID == nil {
			created, next, err := CreateRepeatTodoIfNeed(*todo)
			if err != nil {
				return err
			}

			if created {
				todo.NextID = &next.ID
			}
		}
	}

	plan, err := UpdateTodoRepeatPlan(todo.TodoRepeatPlan, oldTodo.TodoRepeatPlan)
	if err != nil {
		return err
	}
	todo.TodoRepeatPlanID = plan.ID

	// Save
	if err = dal.SaveTodo(todo); err != nil {
		return err
	}

	go UpdateTagAsync(todo, oldTodo.Title)

	return nil
}

func DeleteTodo(userID, todoID int64) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, err
	}

	if err = dal.DeleteTodo(todoID); err != nil {
		return entity.Todo{}, fmt.Errorf("fails to delete todo: %w", err)
	}

	go UpdateTagAsync(&todo, "")

	return todo, nil
}

func OwnTodo(userID, todoID int64) (entity.Todo, error) {
	todo, err := dal.SelectTodo(todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if _, err = OwnOrSharedTodoList(userID, todo.TodoListID); err != nil {
		return entity.Todo{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID)
	}

	return todo, nil
}
