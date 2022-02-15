package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodoRepeatPlan(plan *entity.TodoRepeatPlan) error {
	re := db.Create(plan)
	return utils.WrapGormErr(re.Error, "todo repeat plan")
}

func SelectTodoRepeatPlan(id string) (entity.TodoRepeatPlan, error) {
	var plan entity.TodoRepeatPlan
	re := db.Where("id = ?", id).First(&plan)
	return plan, utils.WrapGormErr(re.Error, "todo repeat plan")
}

func SaveTodoRepeatPlan(todoRepeatPlan *entity.TodoRepeatPlan) error {
	re := db.Save(&todoRepeatPlan)
	return utils.WrapGormErr(re.Error, "todo repeat plan")
}

func DeleteTodoRepeatPlan(id string) error {
	re := db.Delete(&entity.TodoRepeatPlan{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return utils.WrapGormErr(re.Error, "todo repeat plan")
}
