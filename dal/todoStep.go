package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoStep(step *entity.TodoStep) error {
	re := db.Create(step)
	return util.WrapGormErr(re.Error, "todo step")
}

func SelectTodoStep(id string) (entity.TodoStep, error) {
	var step entity.TodoStep
	re := db.Where("id = ?", id).First(&step)
	return step, util.WrapGormErr(re.Error, "todo step")
}

func SelectTodoSteps(todoID string) ([]entity.TodoStep, error) {
	var steps []entity.TodoStep
	re := db.Where(entity.TodoStep{TodoID: todoID}).Find(&steps)
	return steps, util.WrapGormErr(re.Error, "todo step")
}

func SaveTodoStep(todoStep *entity.TodoStep) error {
	re := db.Save(&todoStep)
	return util.WrapGormErr(re.Error, "todo step")
}

func DeleteTodoStep(id string) error {
	re := db.Delete(&entity.TodoStep{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return util.WrapGormErr(re.Error, "todo step")
}
