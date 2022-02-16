package bll

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func CreateTodo(userID string, todo *entity.Todo) error {
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

func GetTodo(userID, todoID string) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %w", err)
	}

	return todo, nil
}

func GetTodos(userID, todoListID string) ([]entity.Todo, error) {
	if _, err := OwnOrSharedTodoList(userID, todoListID); err != nil {
		return nil, err
	}

	todos, err := dal.SelectTodos(todoListID)
	if err != nil {
		return nil, fmt.Errorf("fails to get todos: %w", err)
	}

	return todos, nil
}

func GetImportantTodos(userID string) ([]entity.Todo, error) {
	todos, err := dal.SelectImportantTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get important todos: %w", err)
	}

	return todos, nil
}

func GetPlannedTodos(userID string) ([]entity.Todo, error) {
	todos, err := dal.SelectPlanedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get planed todos: %w", err)
	}

	return todos, nil
}

func GetNotNotifiedTodos(userID string) ([]entity.Todo, error) {
	todos, err := dal.SelectNotNotifiedTodos(userID)
	if err != nil {
		return nil, fmt.Errorf("fails to get not-notified todos: %w", err)
	}

	return todos, nil
}

func UpdateTodo(userID string, todo *entity.Todo) error {
	// Limits
	oldTodo, err := OwnTodo(userID, todo.ID)
	if err != nil {
		return err
	}

	if oldTodo.UserID != todo.UserID {
		return util.NewErrorWithPreconditionFailed("unable to update todo owner")
	}

	// Update values
	if todo.Done && todo.DoneAt.IsZero() {
		t := time.Now()
		todo.DoneAt = &t
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

	// Events
	// TODO[perf] Following events can be async
	if err = UpdateTag(todo, oldTodo.Title); err != nil {
		return err
	}

	if !oldTodo.Done && todo.Done {
		// TODO[feat] Notify new todo
		if _, _, err = CreateRepeatTodoIfNeed(*todo); err != nil {
			return err
		}
	}

	return nil
}

func DeleteTodo(userID, todoID string) (entity.Todo, error) {
	todo, err := OwnTodo(userID, todoID)
	if err != nil {
		return entity.Todo{}, err
	}

	if err = dal.DeleteTodo(todoID); err != nil {
		return entity.Todo{}, fmt.Errorf("fails to delete todo: %v", todoID)
	}

	if err = UpdateTag(&todo, ""); err != nil {
		return entity.Todo{}, err
	}

	return todo, nil
}

func OwnTodo(userID, todoID string) (entity.Todo, error) {
	todo, err := dal.SelectTodo(todoID)
	if err != nil {
		return entity.Todo{}, fmt.Errorf("fails to get todo: %v", todoID)
	}

	if _, err = OwnOrSharedTodoList(userID, todo.TodoListID); err != nil {
		return entity.Todo{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", todo.ID)
	}

	return todo, nil
}
