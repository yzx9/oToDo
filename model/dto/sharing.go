package dto

import "time"

type SharingToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

type SharingTodoList struct {
	TodoListName string `json:"todoListName"`
	UserNickname string `json:"userNickname"`
}
