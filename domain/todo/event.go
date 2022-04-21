package todo

import (
	"encoding/json"
)

func HandleUserCreatedEvent(payload []byte) {
	var dto struct {
		UserID int64 `json:"userID"`
	}

	if err := json.Unmarshal(payload, &dto); err != nil {
		return
	}

	if dto.UserID == 0 {
		// TODO: log, invalid payload, event maybe changed
		return
	}

	if err := TodoListRepository.Save(&TodoList{
		Name:    "Todos", // TODO i18n
		IsBasic: true,
		UserID:  dto.UserID,
	}); err != nil {
		// TODO: log, fails to create user basic todo list
		return
	}

	return
}
