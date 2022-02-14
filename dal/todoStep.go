package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodoStep(step *entity.TodoStep) error {
	re := db.Create(step)
	return utils.WrapGormErr(re.Error, "todo step")
}

func SelectTodoStep(id string) (entity.TodoStep, error) {
	var step entity.TodoStep
	re := db.Where("ID = ?", id).First(&step)
	return step, utils.WrapGormErr(re.Error, "todo step")
}

func SelectTodoSteps(todoID string) ([]entity.TodoStep, error) {
	var steps []entity.TodoStep
	re := db.Where("TodoID = ?", todoID).Find(&steps)
	return steps, utils.WrapGormErr(re.Error, "todo step")
}

func SaveTodoStep(todoStep *entity.TodoStep) error {
	re := db.Save(&todoStep)
	return utils.WrapGormErr(re.Error, "todo step")
}

func DeleteTodoStep(id string) error {
	re := db.Delete(&entity.TodoStep{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return utils.WrapGormErr(re.Error, "todo step")
}
