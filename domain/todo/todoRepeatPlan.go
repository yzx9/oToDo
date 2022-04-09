package todo

import (
	"fmt"
	"time"
)

type TodoRepeatPlan struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Type     string
	Interval int
	Before   *time.Time
	Weekday  int8 // BitBools, [0..6]=[Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday]

	Todos []int64
}

func CreateTodoRepeatPlan(plan TodoRepeatPlan) (TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) {
		return TodoRepeatPlan{}, nil
	}

	plan.ID = 0
	if err := TodoRepeatPlanRepository.Save(&plan); err != nil {
		return TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func UpdateTodoRepeatPlan(plan, oldPlan TodoRepeatPlan) (TodoRepeatPlan, error) {
	if !isValidTodoRepeatPlan(plan) || isSameTodoRepeatPlan(plan, oldPlan) {
		return oldPlan, nil
	}

	plan.ID = 0
	if err := TodoRepeatPlanRepository.Save(&plan); err != nil {
		return TodoRepeatPlan{}, fmt.Errorf("fails to create todo repeat plan: %w", err)
	}

	return plan, nil
}

func GetTodoRepeatPlan(id int64) (TodoRepeatPlan, error) {
	plan, err := TodoRepeatPlanRepository.Find(id)
	if err != nil {
		return TodoRepeatPlan{}, fmt.Errorf("fails to get todo repeat plan: %v", err)
	}

	return plan, nil
}

func CreateRepeatTodoIfNeed(todo Todo) (bool, Todo, error) {
	if todo.TodoRepeatPlanID == 0 {
		return false, Todo{}, nil
	}

	nextDeadline := getTodoNextRepeatTime(todo)
	if todo.TodoRepeatPlan.Before.Before(nextDeadline) {
		return false, Todo{}, nil
	}

	todo.ID = 0
	todo.Deadline = &nextDeadline
	if err := TodoRepository.Save(&todo); err != nil {
		return false, Todo{}, fmt.Errorf("fails to create todo: %w", err)
	}

	return true, todo, nil
}

func isValidTodoRepeatPlan(plan TodoRepeatPlan) bool {
	t := TodoRepeatPlanType(plan.Type)
	if t != TodoRepeatPlanTypeDay &&
		t != TodoRepeatPlanTypeMonth &&
		t != TodoRepeatPlanTypeYear &&
		t != TodoRepeatPlanTypeWeek {
		return false
	}

	// Do not allow set all weekday to false
	if t == TodoRepeatPlanTypeWeek && plan.Weekday == 0 {
		return false
	}

	return plan.Interval > 0
}

func isSameTodoRepeatPlan(plan, oldPlan TodoRepeatPlan) bool {
	if plan.Type != oldPlan.Type ||
		plan.Interval != oldPlan.Interval ||
		plan.Before != oldPlan.Before {
		return false
	}

	if plan.Type == string(TodoRepeatPlanTypeWeek) {
		if plan.Weekday != oldPlan.Weekday {
			return false
		}
	}

	return true
}

func getTodoNextRepeatTime(todo Todo) time.Time {
	deadline := *todo.Deadline
	interval := todo.TodoRepeatPlan.Interval

	weekend := time.Sunday // TODO 此处默认周一为一周开始

	switch TodoRepeatPlanType(todo.TodoRepeatPlan.Type) {
	case TodoRepeatPlanTypeDay:
		return deadline.AddDate(0, 0, interval)

	case TodoRepeatPlanTypeMonth:
		return deadline.AddDate(0, interval, 0)

	case TodoRepeatPlanTypeYear:
		return deadline.AddDate(interval, 0, 0)

	case TodoRepeatPlanTypeWeek:
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
