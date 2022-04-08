package repository

import (
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoStep struct {
	Entity

	Name   string     `json:"name" gorm:"-"`
	Done   bool       `json:"done"`
	DoneAt *time.Time `json:"doneAt"`

	TodoID int64 `json:"todoID"`
	Todo   Todo  `json:"-"`
}

var TodoStepRepo TodoStepRepository

type TodoStepRepository struct {
	db *gorm.DB
}

func (r TodoStepRepository) Save(todoStep *TodoStep) error {
	err := r.db.Save(&todoStep).Error
	return util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) Find(id int64) (TodoStep, error) {
	var step TodoStep
	err := r.db.
		Where(&TodoStep{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&step).
		Error

	return step, util.WrapGormErr(err, "todo step")
}

func (r TodoStepRepository) FindAllByTodo(todoID int64) ([]TodoStep, error) {
	var steps []TodoStep
	err := r.db.
		Where(TodoStep{
			TodoID: todoID,
		}).
		Find(&steps).
		Error

	return steps, util.WrapGormErr(err, "todo step")
}
