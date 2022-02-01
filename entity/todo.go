package entity

import "github.com/google/uuid"

type Todo struct {
	ID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Title string
}
