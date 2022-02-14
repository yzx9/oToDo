package entity

import "time"

type TodoRepeatPlanType string

const (
	TodoRepeatPlanTypeDay   TodoRepeatPlanType = "day"
	TodoRepeatPlanTypeWeek  TodoRepeatPlanType = "week"
	TodoRepeatPlanTypeMonth TodoRepeatPlanType = "month"
	TodoRepeatPlanTypeYear  TodoRepeatPlanType = "year"
)

type TodoRepeatPlan struct {
	ID       string    `json:"-"`
	Type     string    `json:"type"`
	Interval int       `json:"interval"`
	Before   time.Time `json:"before"`
	Weekday  [7]byte   `json:"weekday"` // Follow time.Weekend: Sunday Monday Tuesday Wednesday Thursday Friday Saturday
}
