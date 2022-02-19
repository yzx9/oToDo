package dto

import "time"

type SharingToken struct {
	Token     string    `json:"token"`
	Type      int8      `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type SharingTodoList struct {
	TodoListName string `json:"todoListName"`
	UserNickname string `json:"userNickname"`
}
