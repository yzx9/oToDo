package entity

import (
	"time"
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

	Type     string    `json:"type" gorm:"size:8"`
	Interval int       `json:"interval"`
	Before   time.Time `json:"before"`
	Weekday  int8      `json:"weekday"` // BitBools, [0..6]=[Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday]

	Todos []Todo `json:"-"`
}
