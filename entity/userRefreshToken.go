package entity

import "time"

type UserRefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenID   string    `json:"token_id"`
	CreatedAt time.Time `json:"created_at"`
}
