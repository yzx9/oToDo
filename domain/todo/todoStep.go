package todo

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

type TodoStep struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name   string
	Done   bool
	DoneAt *time.Time

	TodoID int64
}

func CreateTodoStep(userID, todoID int64, name string) (TodoStep, error) {
	_, err := OwnTodo(userID, todoID)
	if err != nil {
		return TodoStep{}, util.NewErrorWithNotFound("todo not found: %v", todoID)
	}

	step := TodoStep{
		Name:   name,
		TodoID: todoID,
	}
	if err = TodoStepRepository.Save(&step); err != nil {
		return TodoStep{}, util.NewErrorWithUnknown("fails to create todo step")
	}

	return step, nil
}

func UpdateTodoStep(userID int64, step *TodoStep) error {
	oldStep, err := OwnTodoStep(userID, step.ID)
	if err != nil {
		return err
	}

	step.CreatedAt = oldStep.CreatedAt
	step.TodoID = oldStep.TodoID

	if step.Done && !oldStep.Done {
		t := time.Now()
		step.DoneAt = &t
	}

	if err = TodoStepRepository.Save(step); err != nil {
		return util.NewErrorWithUnknown("fails to update todo step")
	}

	return nil
}

func DeleteTodoStep(userID, todoID, todoStepID int64) (TodoStep, error) {
	step, err := OwnTodoStep(userID, todoStepID)
	if err != nil {
		return TodoStep{}, err
	}

	if step.TodoID != todoID {
		return TodoStep{}, util.NewErrorWithNotFound("todo step not found in todo: %v", todoStepID)
	}

	return step, TodoStepRepository.Delete(todoStepID)
}

func OwnTodoStep(userID, todoStepID int64) (TodoStep, error) {
	step, err := TodoStepRepository.Find(todoStepID)
	if err != nil {
		return TodoStep{}, fmt.Errorf("fails to get todo step: %w", err)
	}

	_, err = OwnTodo(userID, step.TodoID)
	if err != nil {
		return TodoStep{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", step.ID)
	}

	return step, nil
}
