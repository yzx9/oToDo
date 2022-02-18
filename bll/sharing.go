package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
)

func CreateTodoListSharing(userID, todoListID string) (entity.Sharing, error) {
	exist, err := dal.ExistActiveSharing(userID, entity.SharingTypeTodoList)
	if err != nil {
		return entity.Sharing{}, fmt.Errorf("fails to check todo list sharing: %w", err)
	}

	// inactive old sharings
	if exist {
		sharings, err := dal.SelectActiveSharings(userID, entity.SharingTypeTodoList)
		if err != nil {
			return entity.Sharing{}, fmt.Errorf("fails to check todo list sharing: %w", err)
		}

		for i := range sharings {
			sharings[i].Active = false
			if err := dal.SaveSharing(&sharings[i]); err != nil {
				return entity.Sharing{}, fmt.Errorf("fails to inactive todo list sharing: %w", err)
			}
		}
	}

	sharing := entity.Sharing{
		Token:   "TODO", // TODO
		Active:  true,
		Type:    entity.SharingTypeTodoList,
		Payload: "TODO", // TODO
		UserID:  userID,
	}
	if err := dal.InsertSharing(&sharing); err != nil {
		return entity.Sharing{}, fmt.Errorf("fails to create todo list sharing: %w", err)
	}

	return sharing, nil
}
