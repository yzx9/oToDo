package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func CreateTodoRepeatPlan(plan entity.TodoRepeatPlan) (entity.TodoRepeatPlan, error) {
	if !IsValidTodoRepeatPlan(plan) {
		return entity.TodoRepeatPlan{
			ID: "",
		}, nil
	}

	plan.ID = uuid.NewString()
	if err := dal.InsertTodoRepeatPlan(plan); err != nil {
		return entity.TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func UpdateTodoRepeatPlan(plan, oldPlan entity.TodoRepeatPlan) (entity.TodoRepeatPlan, error) {
	if !IsValidTodoRepeatPlan(plan) {
		return entity.TodoRepeatPlan{
			ID: "",
		}, nil
	}

	if !IsSameTodoRepeatPlan(plan, oldPlan) {
		return oldPlan, nil
	}

	plan.ID = uuid.NewString()
	if err := dal.InsertTodoRepeatPlan(plan); err != nil {
		return entity.TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func GetTodoRepeatPlan(id string) (entity.TodoRepeatPlan, error) {
	plan, err := dal.GetTodoRepeatPlan(id)
	if err != nil {
		return entity.TodoRepeatPlan{}, fmt.Errorf("fails to get todo repeat plan: %v", err)
	}

	return plan, nil
}

func IsValidTodoRepeatPlan(plan entity.TodoRepeatPlan) bool {
	return plan.Interval != 0
}

func IsSameTodoRepeatPlan(plan, oldPlan entity.TodoRepeatPlan) bool {
	return plan.Before == oldPlan.Before && plan.Interval == oldPlan.Interval
}

func CreateRepeatTodoIfNeed(todo entity.Todo) (bool, entity.Todo, error) {
	duration := time.Duration(int64(time.Second) * int64(todo.RepeatPlan.Interval))
	nextDeadline := todo.Deadline.Add(duration)
	if todo.RepeatPlanID == "" || todo.RepeatPlan.Before.Before(nextDeadline) {
		return false, entity.Todo{}, nil
	}

	todo.ID = uuid.NewString()
	todo.Deadline = nextDeadline
	if err := dal.InsertTodo(todo); err != nil {
		return false, entity.Todo{}, fmt.Errorf("fails to create todo: %w", err)
	}

	return true, todo, nil
}
