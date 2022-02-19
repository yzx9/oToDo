package dto

import "time"

type Sharing struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}
