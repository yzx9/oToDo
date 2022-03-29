package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoRepeatPlan(plan *entity.TodoRepeatPlan) error {
	err := db.Create(plan).Error
	return util.WrapGormErr(err, "todo repeat plan")
}

func SelectTodoRepeatPlan(id int64) (entity.TodoRepeatPlan, error) {
	var plan entity.TodoRepeatPlan
	err := db.
		Where(&entity.TodoRepeatPlan{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&plan).
		Error

	return plan, util.WrapGormErr(err, "todo repeat plan")
}

func SaveTodoRepeatPlan(todoRepeatPlan *entity.TodoRepeatPlan) error {
	err := db.Save(&todoRepeatPlan).Error
	return util.WrapGormErr(err, "todo repeat plan")
}

func DeleteTodoRepeatPlan(id int64) error {
	err := db.
		Delete(&entity.TodoRepeatPlan{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo repeat plan")
}
