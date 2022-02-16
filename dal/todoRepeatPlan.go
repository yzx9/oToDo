package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoRepeatPlan(plan *entity.TodoRepeatPlan) error {
	re := db.Create(plan)
	return util.WrapGormErr(re.Error, "todo repeat plan")
}

func SelectTodoRepeatPlan(id int64) (entity.TodoRepeatPlan, error) {
	var plan entity.TodoRepeatPlan
	re := db.Where("id = ?", id).First(&plan)
	return plan, util.WrapGormErr(re.Error, "todo repeat plan")
}

func SaveTodoRepeatPlan(todoRepeatPlan *entity.TodoRepeatPlan) error {
	re := db.Save(&todoRepeatPlan)
	return util.WrapGormErr(re.Error, "todo repeat plan")
}

func DeleteTodoRepeatPlan(id int64) error {
	re := db.Delete(&entity.TodoRepeatPlan{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return util.WrapGormErr(re.Error, "todo repeat plan")
}
