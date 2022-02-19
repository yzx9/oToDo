package bll

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

/**
 * oTodo List Sharing
 */
func CreateTodoListSharing(userID, todoListID int64) (entity.Sharing, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return entity.Sharing{}, err
	}

	if todoList.IsBasic {
		return entity.Sharing{}, fmt.Errorf("unable to share basic todo list: %v", todoListID)
	}

	// Only allow one sharing active
	if _, err = dal.DeleteSharings(userID, entity.SharingTypeTodoList); err != nil {
		return entity.Sharing{}, fmt.Errorf("fails to delete old sharing tokens: %w", err)
	}

	sharing := entity.Sharing{
		Token:     newSharingToken(),
		Active:    true,
		Type:      entity.SharingTypeTodoList,
		RelatedID: todoListID,
		UserID:    userID,
	}
	if err := dal.InsertSharing(&sharing); err != nil {
		return entity.Sharing{}, fmt.Errorf("fails to create sharing token: %w", err)
	}

	return sharing, nil
}

func GetActiveTodoListSharings(userID, todoListID int64) ([]entity.Sharing, error) {
	sharings, err := dal.SelectActiveSharings(userID, entity.SharingTypeTodoList)
	if err != nil {
		return nil, fmt.Errorf("fails to get sharing tokens: %w", err)
	}

	vec := make([]entity.Sharing, 0)
	for i := range sharings {
		if sharings[i].RelatedID == todoListID {
			vec = append(vec, sharings[i])
		}
	}

	return vec, nil
}

func DeleteTodoListSharing(userID int64, token string) error {
	sharing, err := ValidSharing(token)
	if err != nil {
		return err
	}

	if sharing.Type != entity.SharingTypeTodoList {
		return util.NewErrorWithForbidden("invalid sharing token: %v")
	}

	if sharing.UserID != userID {
		return util.NewErrorWithForbidden("unable to delete non-own sharing token")
	}

	sharing.Active = false
	if err := dal.SaveSharing(&sharing); err != nil {
		return fmt.Errorf("fails to delete sharing: %w", err)
	}

	return nil
}

func ValidSharing(token string) (entity.Sharing, error) {
	sharing, err := dal.SelectSharing(token)
	if err != nil {
		return entity.Sharing{}, fmt.Errorf("invalid sharing token: %w", err)
	}

	if !sharing.Active {
		return entity.Sharing{}, util.NewErrorWithForbidden("sharing token has been inactive: %v", token)
	}

	return sharing, nil
}

/**
 * common
 */

func newSharingToken() string {
	return base64.RawStdEncoding.EncodeToString([]byte(uuid.NewString()))
}
