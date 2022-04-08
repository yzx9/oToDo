package todo

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

func CreateTodo(userID int64, todo *repository.Todo) error {
	_, err := todolist.OwnOrSharedTodoList(userID, todo.TodoListID)
	if err != nil {
		return fmt.Errorf("fails to get todo list: %w", err)
	}

	todo.UserID = userID // override user

	plan, err := CreateTodoRepeatPlan(todo.TodoRepeatPlan)
	if err != nil {
		return fmt.Errorf("fails to create todo repeat plan: %w", err)
	}
	todo.TodoRepeatPlanID = plan.ID

	todo.ID = 0
	if err := TodoRepository.Save(todo); err != nil {
		return fmt.Errorf("fails to create todo: %w", err)
	}

	return nil
}

func UpdateTodo(userID int64, todo *repository.Todo) error {
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
	if err = TodoRepository.Save(todo); err != nil {
		return err
	}

	go UpdateTagAsync(todo, oldTodo.Title)

	return nil
}

func DeleteTodo(userID, todoID int64) (repository.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return repository.Todo{}, err
	}

	if err = TodoRepository.Delete(todoID); err != nil {
		return repository.Todo{}, fmt.Errorf("fails to delete todo: %w", err)
	}

	go UpdateTagAsync(&todo, "")

	return todo, nil
}

func OwnTodo(userID, todoID int64) (repository.Todo, error) {
	todo, err := TodoRepository.Find(todoID)
	if err != nil {
		return repository.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	if _, err = todolist.OwnOrSharedTodoList(userID, todo.TodoListID); err != nil {
		return repository.Todo{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID)
	}

	return todo, nil
}
