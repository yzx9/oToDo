package repository

import (
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoRepeatPlanType string

const (
	TodoRepeatPlanTypeDay   TodoRepeatPlanType = "day"
	TodoRepeatPlanTypeWeek  TodoRepeatPlanType = "week"
	TodoRepeatPlanTypeMonth TodoRepeatPlanType = "month"
	TodoRepeatPlanTypeYear  TodoRepeatPlanType = "year"
)

type TodoRepeatPlan struct {
	Entity

	Type     string     `json:"type" gorm:"size:8"`
	Interval int        `json:"interval"`
	Before   *time.Time `json:"before"`
	Weekday  int8       `json:"weekday"` // BitBools, [0..6]=[Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday]

	Todos []Todo `json:"-"`
}

var TodoRepeatPlanRepo TodoRepeatPlanRepository

type TodoRepeatPlanRepository struct {
	db *gorm.DB
}

func NewTodoRepeatPlanRepository(db *gorm.DB) TodoRepeatPlanRepository {
	return TodoRepeatPlanRepository{db: db}
}

func (r TodoRepeatPlanRepository) Save(plan *TodoRepeatPlan) error {
	err := r.db.Save(plan).Error
	return util.WrapGormErr(err, "todo repeat plan")
}

func (r TodoRepeatPlanRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoRepeatPlan{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo repeat plan")
}

func (r TodoRepeatPlanRepository) Find(id int64) (TodoRepeatPlan, error) {
	var plan TodoRepeatPlan
	err := r.db.
		Where(&TodoRepeatPlan{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&plan).
		Error

	return plan, util.WrapGormErr(err, "todo repeat plan")
}
