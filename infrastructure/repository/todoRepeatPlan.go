package repository

import (
	"time"

	"github.com/yzx9/otodo/infrastructure/util"
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

func InsertTodoRepeatPlan(plan *TodoRepeatPlan) error {
	err := db.Create(plan).Error
	return util.WrapGormErr(err, "todo repeat plan")
}

func SelectTodoRepeatPlan(id int64) (TodoRepeatPlan, error) {
	var plan TodoRepeatPlan
	err := db.
		Where(&TodoRepeatPlan{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&plan).
		Error

	return plan, util.WrapGormErr(err, "todo repeat plan")
}

func SaveTodoRepeatPlan(todoRepeatPlan *TodoRepeatPlan) error {
	err := db.Save(&todoRepeatPlan).Error
	return util.WrapGormErr(err, "todo repeat plan")
}

func DeleteTodoRepeatPlan(id int64) error {
	err := db.
		Delete(&TodoRepeatPlan{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo repeat plan")
}
