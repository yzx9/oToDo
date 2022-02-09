package entity

import (
	"time"
)

type User struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Nickname        string    `json:"nickname"`
	Password        []byte    `json:"password"`
	Email           string    `json:"email"`
	Telephone       string    `json:"telephone"`
	Avatar          string    `json:"avatar"`
	BasicTodoListID string    `json:"basic_todo_list_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserRefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenID   string    `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
}
