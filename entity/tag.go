package entity

import "github.com/google/uuid"

type Tag struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title      string
	CountTotal uint
	CountTodo  uint
	CountDone  uint

	UserId uuid.UUID
	User   User
}
