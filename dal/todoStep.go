package dal

import (
	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var todoSteps = make(map[uuid.UUID]entity.TodoStep)

func InsertTodoStep(step entity.TodoStep) error {
	todoSteps[step.ID] = step
	return nil
}

func GetTodoStep(id uuid.UUID) (entity.TodoStep, error) {
	step, ok := todoSteps[id]
	if !ok {
		return entity.TodoStep{}, utils.NewErrorWithNotFound("todo step not found: %v", id)
	}

	return step, nil
}

func GetTodoSteps(todoID uuid.UUID) ([]entity.TodoStep, error) {
	vec := make([]entity.TodoStep, 0)
	for _, v := range todoSteps {
		if v.TodoID == todoID {
			vec = append(vec, v)
		}
	}

	return vec, nil
}

func UpdateTodoStep(todoStep entity.TodoStep) error {
	_, exists := todoSteps[todoStep.ID]
	if !exists {
		return utils.NewErrorWithNotFound("todo step not found: %v", todoStep.ID)
	}

	todoSteps[todoStep.ID] = todoStep
	return nil
}

func DeleteTodoStep(id uuid.UUID) error {
	_, exists := todoSteps[id]
	if !exists {
		return utils.NewErrorWithNotFound("todo step not found: %v", id)
	}

	delete(todoSteps, id)
	return nil
}
