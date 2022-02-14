package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func CreateTodoRepeatPlan(plan entity.TodoRepeatPlan) (entity.TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) {
		return entity.TodoRepeatPlan{
			Entity: entity.Entity{
				ID: "",
			},
		}, nil
	}

	plan.ID = uuid.NewString()
	if err := dal.InsertTodoRepeatPlan(plan); err != nil {
		return entity.TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func UpdateTodoRepeatPlan(plan, oldPlan entity.TodoRepeatPlan) (entity.TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) {
		return entity.TodoRepeatPlan{
			Entity: entity.Entity{
				ID: "",
			},
		}, nil
	}

	if !isSameTodoRepeatPlan(plan, oldPlan) {
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

func CreateRepeatTodoIfNeed(todo entity.Todo) (bool, entity.Todo, error) {
	if todo.TodoRepeatPlanID == "" {
		return false, entity.Todo{}, nil
	}

	nextDeadline := getTodoNextRepeatTime(todo)
	if todo.TodoRepeatPlan.Before.Before(nextDeadline) {
		return false, entity.Todo{}, nil
	}

	todo.ID = uuid.NewString()
	todo.Deadline = nextDeadline
	if err := dal.InsertTodo(todo); err != nil {
		return false, entity.Todo{}, fmt.Errorf("fails to create todo: %w", err)
	}

	return true, todo, nil
}

func isValidTodoRepeatPlan(plan entity.TodoRepeatPlan) bool {
	t := entity.TodoRepeatPlanType(plan.Type)
	if t != entity.TodoRepeatPlanTypeDay &&
		t != entity.TodoRepeatPlanTypeMonth &&
		t != entity.TodoRepeatPlanTypeYear &&
		t != entity.TodoRepeatPlanTypeWeek {
		return false
	}

	// Do not allow set all weekday to false
	if t == entity.TodoRepeatPlanTypeWeek && plan.Weekday == 0 {
		return false
	}

	return plan.Interval > 0
}

func isSameTodoRepeatPlan(plan, oldPlan entity.TodoRepeatPlan) bool {
	if plan.Type != oldPlan.Type ||
		plan.Interval != oldPlan.Interval ||
		plan.Before != oldPlan.Before {
		return false
	}

	if plan.Type == string(entity.TodoRepeatPlanTypeWeek) {
		if plan.Weekday != oldPlan.Weekday {
			return false
		}
	}

	return true
}

func getTodoNextRepeatTime(todo entity.Todo) time.Time {
	deadline := todo.Deadline
	interval := todo.TodoRepeatPlan.Interval

	weekend := time.Sunday // TODO 此处默认周一为一周开始

	switch entity.TodoRepeatPlanType(todo.TodoRepeatPlan.Type) {
	case entity.TodoRepeatPlanTypeDay:
		return deadline.AddDate(0, 0, interval)

	case entity.TodoRepeatPlanTypeMonth:
		return deadline.AddDate(0, interval, 0)

	case entity.TodoRepeatPlanTypeYear:
		return deadline.AddDate(interval, 0, 0)

	case entity.TodoRepeatPlanTypeWeek:
		if deadline.Weekday() == weekend {
			deadline = deadline.AddDate(0, 0, (interval-1)*7)
		}
		deadline = deadline.AddDate(0, 0, 1)
		for i := 0; i < 7-1; i++ {
			mask := int8(0x01 << deadline.Weekday())
			if todo.TodoRepeatPlan.Weekday&mask != 0 {
				return deadline
			}
			deadline = deadline.AddDate(0, 0, 1)
		}
		return time.Time{}

	default:
		return time.Time{}
	}
}
