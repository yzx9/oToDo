package bll

import (
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

/**
 * oTodo List Sharing
 */
func CreateTodoListSharing(userID, todoListID int64) (repository.Sharing, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return repository.Sharing{}, err
	}

	if todoList.IsBasic {
		return repository.Sharing{}, fmt.Errorf("unable to share basic todo list: %v", todoListID)
	}

	// Only allow one sharing active
	if _, err = repository.SharingRepo.DeleteSharings(userID, repository.SharingTypeTodoList); err != nil {
		return repository.Sharing{}, fmt.Errorf("fails to delete old sharing tokens: %w", err)
	}

	sharing := repository.Sharing{
		Token:     newSharingToken(),
		Active:    true,
		Type:      repository.SharingTypeTodoList,
		RelatedID: todoListID,
		UserID:    userID,
	}
	if err := repository.SharingRepo.InsertSharing(&sharing); err != nil {
		return repository.Sharing{}, fmt.Errorf("fails to create sharing token: %w", err)
	}

	return sharing, nil
}

func GetActiveTodoListSharings(userID, todoListID int64) ([]repository.Sharing, error) {
	sharings, err := repository.SharingRepo.SelectActiveSharings(userID, repository.SharingTypeTodoList)
	if err != nil {
		return nil, fmt.Errorf("fails to get sharing tokens: %w", err)
	}

	vec := make([]repository.Sharing, 0)
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

	if sharing.Type != repository.SharingTypeTodoList {
		return util.NewErrorWithForbidden("invalid sharing token: %v")
	}

	if sharing.UserID != userID {
		return util.NewErrorWithForbidden("unable to delete non-own sharing token")
	}

	sharing.Active = false
	if err := repository.SharingRepo.SaveSharing(&sharing); err != nil {
		return fmt.Errorf("fails to delete sharing: %w", err)
	}

	return nil
}

func ValidSharing(token string) (repository.Sharing, error) {
	sharing, err := repository.SharingRepo.SelectSharing(token)
	if err != nil {
		return repository.Sharing{}, fmt.Errorf("invalid sharing token: %w", err)
	}

	if !sharing.Active {
		return repository.Sharing{}, util.NewErrorWithForbidden("sharing token has been inactive: %v", token)
	}

	return sharing, nil
}

/**
 * common
 */

func newSharingToken() string {
	return base64.RawStdEncoding.EncodeToString([]byte(uuid.NewString()))
}
