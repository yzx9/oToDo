package repository

import (
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
)

type TodoStep struct {
	Entity

	Name   string     `json:"name" gorm:"-"`
	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	TodoID int64 `json:"todoID"`
	Todo   Todo  `json:"-"`
}

func InsertTodoStep(step *TodoStep) error {
	err := db.Create(step).Error
	return util.WrapGormErr(err, "todo step")
}

func SelectTodoStep(id int64) (TodoStep, error) {
	var step TodoStep
	err := db.
		Where(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&step).
		Error

	return step, util.WrapGormErr(err, "todo step")
}

func SelectTodoSteps(todoID int64) ([]TodoStep, error) {
	var steps []TodoStep
	err := db.
		Where(TodoStep{
			TodoID: todoID,
		}).
		Find(&steps).
		Error

	return steps, util.WrapGormErr(err, "todo step")
}

func SaveTodoStep(todoStep *TodoStep) error {
	err := db.Save(&todoStep).Error
	return util.WrapGormErr(err, "todo step")
}

func DeleteTodoStep(id int64) error {
	err := db.
		Delete(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo step")
}
