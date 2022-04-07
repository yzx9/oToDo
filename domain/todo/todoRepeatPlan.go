package todo

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/repository"
)

func CreateTodoRepeatPlan(plan repository.TodoRepeatPlan) (repository.TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) {
		return repository.TodoRepeatPlan{}, nil
	}

	plan.ID = 0
	if err := repository.TodoRepeatPlanRepo.Save(&plan); err != nil {
		return repository.TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func UpdateTodoRepeatPlan(plan, oldPlan repository.TodoRepeatPlan) (repository.TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) || isSameTodoRepeatPlan(plan, oldPlan) {
		return oldPlan, nil
	}

	plan.ID = 0
	if err := repository.TodoRepeatPlanRepo.Save(&plan); err != nil {
		return repository.TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func GetTodoRepeatPlan(id int64) (repository.TodoRepeatPlan, error) {
	plan, err := repository.TodoRepeatPlanRepo.Find(id)
	if err != nil {
		return repository.TodoRepeatPlan{}, fmt.Errorf("fails to get todo repeat plan: %v", err)
	}

	return plan, nil
}

func CreateRepeatTodoIfNeed(todo repository.Todo) (bool, repository.Todo, error) {
	if todo.TodoRepeatPlanID == 0 {
		return false, repository.Todo{}, nil
	}

	nextDeadline := getTodoNextRepeatTime(todo)
	if todo.TodoRepeatPlan.Before.Before(nextDeadline) {
		return false, repository.Todo{}, nil
	}

	todo.ID = 0
	todo.Deadline = &nextDeadline
	if err := repository.TodoRepo.Save(&todo); err != nil {
		return false, repository.Todo{}, fmt.Errorf("fails to create todo: %w", err)
	}

	return true, todo, nil
}

func isValidTodoRepeatPlan(plan repository.TodoRepeatPlan) bool {
	t := repository.TodoRepeatPlanType(plan.Type)
	if t != repository.TodoRepeatPlanTypeDay &&
		t != repository.TodoRepeatPlanTypeMonth &&
		t != repository.TodoRepeatPlanTypeYear &&
		t != repository.TodoRepeatPlanTypeWeek {
		return false
	}

	// Do not allow set all weekday to false
	if t == repository.TodoRepeatPlanTypeWeek && plan.Weekday == 0 {
		return false
	}

	return plan.Interval > 0
}

func isSameTodoRepeatPlan(plan, oldPlan repository.TodoRepeatPlan) bool {
	if plan.Type != oldPlan.Type ||
		plan.Interval != oldPlan.Interval ||
		plan.Before != oldPlan.Before {
		return false
	}

	if plan.Type == string(repository.TodoRepeatPlanTypeWeek) {
		if plan.Weekday != oldPlan.Weekday {
			return false
		}
	}

	return true
}

func getTodoNextRepeatTime(todo repository.Todo) time.Time {
	deadline := *todo.Deadline
	interval := todo.TodoRepeatPlan.Interval

	weekend := time.Sunday // TODO 此处默认周一为一周开始

	switch repository.TodoRepeatPlanType(todo.TodoRepeatPlan.Type) {
	case repository.TodoRepeatPlanTypeDay:
		return deadline.AddDate(0, 0, interval)

	case repository.TodoRepeatPlanTypeMonth:
		return deadline.AddDate(0, interval, 0)

	case repository.TodoRepeatPlanTypeYear:
		return deadline.AddDate(interval, 0, 0)

	case repository.TodoRepeatPlanTypeWeek:
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
