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

func (step TodoStep) New() error {
	if err := TodoStepRepository.Save(&step); err != nil {
		return util.NewErrorWithUnknown("fails to create todo step")
	}

	return nil
}

func (step *TodoStep) Save(userID int64) error {
	oldStep, err := TodoStepRepository.Find(step.ID)
	if err != nil {
		return fmt.Errorf("fails to get todo step: %w", err)
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

func (step *TodoStep) Delete() error {
	// TODO: delte association
	return TodoStepRepository.Delete(step.ID)
}
