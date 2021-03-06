package repository

import (
	"time"

	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type TodoRepeatPlan struct {
	repository.Entity

	Type     string `gorm:"size:8"`
	Interval int
	Before   *time.Time
	Weekday  int8 // BitBools, [0..6]=[Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday]

	Todos []Todo
}

type TodoRepeatPlanRepository struct {
	db *gorm.DB
}

func NewTodoRepeatPlanRepository(db *gorm.DB) TodoRepeatPlanRepository {
	return TodoRepeatPlanRepository{db: db}
}

func (r TodoRepeatPlanRepository) Save(entity *todo.TodoRepeatPlan) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "todo repeat plan")
}

func (r TodoRepeatPlanRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoRepeatPlan{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo repeat plan")
}

func (r TodoRepeatPlanRepository) Find(id int64) (todo.TodoRepeatPlan, error) {
	var po TodoRepeatPlan
	err := r.db.
		Where(&TodoRepeatPlan{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "todo repeat plan")
}

func (r TodoRepeatPlanRepository) convertToPO(entity *todo.TodoRepeatPlan) TodoRepeatPlan {
	return TodoRepeatPlan{
		Entity: repository.Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Type:     entity.Type,
		Interval: entity.Interval,
		Before:   entity.Before,
		Weekday:  entity.Weekday,

		Todos: nil, // TODO
	}
}

func (r TodoRepeatPlanRepository) convertToEntity(po TodoRepeatPlan) todo.TodoRepeatPlan {
	return todo.TodoRepeatPlan{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Type:     po.Type,
		Interval: po.Interval,
		Before:   po.Before,
		Weekday:  po.Weekday,

		Todos: nil, // TODO
	}
}
