package repository

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoStep(step *entity.TodoStep) error {
	err := db.Create(step).Error
	return util.WrapGormErr(err, "todo step")
}

func SelectTodoStep(id int64) (entity.TodoStep, error) {
	var step entity.TodoStep
	err := db.
		Where(&entity.TodoStep{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&step).
		Error

	return step, util.WrapGormErr(err, "todo step")
}

func SelectTodoSteps(todoID int64) ([]entity.TodoStep, error) {
	var steps []entity.TodoStep
	err := db.
		Where(entity.TodoStep{
			TodoID: todoID,
		}).
		Find(&steps).
		Error

	return steps, util.WrapGormErr(err, "todo step")
}

func SaveTodoStep(todoStep *entity.TodoStep) error {
	err := db.Save(&todoStep).Error
	return util.WrapGormErr(err, "todo step")
}

func DeleteTodoStep(id int64) error {
	err := db.
		Delete(&entity.TodoStep{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo step")
}
