package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func CreateTodo(userID string, todo entity.Todo) (entity.Todo, error) {
	todo.ID = uuid.NewString()
	todo.UserID = userID // override user

	plan, err := CreateTodoRepeatPlan(todo.RepeatPlan)
	if err != nil {
		return entity.Todo{}, err
	}
	todo.RepeatPlanID = plan.ID

	if err := dal.InsertTodo(todo); err != nil {
		return entity.Todo{}, fmt.Errorf("fails to create todo: %w", err)
	}

	return todo, nil
}

func GetTodo(userID, todoID string) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	plan, err := GetTodoRepeatPlan(todo.RepeatPlanID)
	if err != nil {
		return entity.Todo{}, err
	}

	todo.RepeatPlan = plan
	return todo, nil
}

func GetTodos(userID, todoListID string) ([]entity.Todo, error) {
	if _, err := OwnTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	todos, err := dal.GetTodos(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodos(userID string) ([]entity.Todo, error) {
	todos, err := dal.GetPlanedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodos(userID string) ([]entity.Todo, error) {
	todos, err := dal.GetNotNotifiedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}

func UpdateTodo(userID string, todo entity.Todo) (entity.Todo, error) {
	// Limits
	oldTodo, err := OwnTodo(userID, todo.ID)
	if err != nil {
		return entity.Todo{}, err
	}

	if oldTodo.UserID != todo.UserID {
		return entity.Todo{}, utils.NewErrorWithPreconditionFailed("unable to update todo owner")
	}

	// Update values
	if !oldTodo.Done && todo.Done {
		todo.DoneAt = time.Now()
	}

	plan, err := UpdateTodoRepeatPlan(todo.RepeatPlan, oldTodo.RepeatPlan)
	if err != nil {
		return entity.Todo{}, err
	}
	todo.RepeatPlanID = plan.ID

	// Save
	if err = dal.UpdateTodo(todo); err != nil {
		return entity.Todo{}, err
	}

	// Update tags
	// TODO How to update shared user
	// TODO Async
	// TODO Catch following error, but dont throw
	UpdateTag(userID, todo.ID, todo.Title, oldTodo.Title)

	// Create new todo if repeat
	if !oldTodo.Done && todo.Done {
		// TODO Async
		// TODO Catch following error, but dont throw
		CreateRepeatTodoIfNeed(todo)
	}

	return todo, nil
}

func DeleteTodo(userID, todoID string) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, err
	}

	err = dal.DeleteTodo(todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to delete todo: %v", todoID)
	}

	// TODO How to update shared user
	// TODO record following error, but dont throw
	UpdateTag(userID, todo.ID, todo.Title, "")

	return todo, nil
}

func OwnTodo(userID, todoID string) (entity.Todo, error) {
	r := func(err error) (entity.Todo, error) {
		return entity.Todo{}, err
	}

	todo, err := dal.GetTodo(todoID)
	if err != nil {
		return r(fmt.Errorf("fails to get todo: %v", todoID))
	}

	// TODO todo in shared todo list
	if todo.UserID != userID {
		return r(utils.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID))
	}

	return todo, nil
}
