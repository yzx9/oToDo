package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoRepeatPlans = make(map[string]entity.TodoRepeatPlan)

func InsertTodoRepeatPlan(plan entity.TodoRepeatPlan) error {
	todoRepeatPlans[plan.ID] = plan
	return nil
}

func GetTodoRepeatPlan(id string) (entity.TodoRepeatPlan, error) {
	plan, ok := todoRepeatPlans[id]
	if !ok {
		return entity.TodoRepeatPlan{}, utils.NewErrorWithNotFound("todo repeat plan not found: %v", id)
	}

	return plan, nil
}

func UpdateTodoRepeatPlan(todoRepeatPlan entity.TodoRepeatPlan) error {
	_, exists := todoRepeatPlans[todoRepeatPlan.ID]
	if !exists {
		return utils.NewErrorWithNotFound("todo repeat plan not found: %v", todoRepeatPlan.ID)
	}

	todoRepeatPlans[todoRepeatPlan.ID] = todoRepeatPlan
	return nil
}

func DeleteTodoRepeatPlan(id string) error {
	_, exists := todoRepeatPlans[id]
	if !exists {
		return utils.NewErrorWithNotFound("todo repeat plan not found: %v", id)
	}

	delete(todoRepeatPlans, id)
	return nil
}
