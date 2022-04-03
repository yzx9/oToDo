package bll

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

func CreateTodoStep(userID, todoID int64, name string) (repository.TodoStep, error) {
	_, err := OwnTodo(userID, todoID)
	if err != nil {
		return repository.TodoStep{}, util.NewErrorWithNotFound("todo not found: %v", todoID)
	}

	step := repository.TodoStep{
		Name:   name,
		TodoID: todoID,
	}
	if err = repository.InsertTodoStep(&step); err != nil {
		return repository.TodoStep{}, util.NewErrorWithUnknown("fails to create todo step")
	}

	return step, nil
}

func UpdateTodoStep(userID int64, step *repository.TodoStep) error {
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

	if err = repository.SaveTodoStep(step); err != nil {
		return util.NewErrorWithUnknown("fails to update todo step")
	}

	return nil
}

func DeleteTodoStep(userID, todoID, todoStepID int64) (repository.TodoStep, error) {
	step, err := OwnTodoStep(userID, todoStepID)
	if err != nil {
		return repository.TodoStep{}, err
	}

	if step.TodoID != todoID {
		return repository.TodoStep{}, util.NewErrorWithNotFound("todo step not found in todo: %v", todoStepID)
	}

	return step, repository.DeleteTodoStep(todoStepID)
}

func OwnTodoStep(userID, todoStepID int64) (repository.TodoStep, error) {
	step, err := repository.SelectTodoStep(todoStepID)
	if err != nil {
		return repository.TodoStep{}, fmt.Errorf("fails to get todo step: %w", err)
	}

	_, err = OwnTodo(userID, step.TodoID)
	if err != nil {
		return repository.TodoStep{}, util.NewErrorWithForbidden("unable to handle non-owned todo: %v", step.ID)
	}

	return step, nil
}
