package todolist

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/infrastructure/util"
)

type Sharing struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Token     string
	Active    bool
	Type      int8  // SharingType
	RelatedID int64 // Depends on Type

	UserID int64
}

func CreateTodoListSharing(userID, todoListID int64) (Sharing, error) {
	todoList, err := OwnTodoList(userID, todoListID)
	if err != nil {
		return Sharing{}, err
	}

	if todoList.IsBasic {
		return Sharing{}, fmt.Errorf("unable to share basic todo list: %v", todoListID)
	}

	// Only allow one sharing active
	if _, err = SharingRepository.DeleteAllByUserAndType(userID, SharingTypeTodoList); err != nil {
		return Sharing{}, fmt.Errorf("fails to delete old sharing tokens: %w", err)
	}

	sharing := Sharing{
		Token:     newSharingToken(),
		Active:    true,
		Type:      SharingTypeTodoList,
		RelatedID: todoListID,
		UserID:    userID,
	}
	if err := SharingRepository.Save(&sharing); err != nil {
		return Sharing{}, fmt.Errorf("fails to create sharing token: %w", err)
	}

	return sharing, nil
}

func DeleteTodoListSharing(userID int64, token string) error {
	sharing, err := GetSharing(token)
	if err != nil {
		return err
	}

	if sharing.Type != SharingTypeTodoList {
		return util.NewErrorWithForbidden("invalid sharing token: %v")
	}

	if sharing.UserID != userID {
		return util.NewErrorWithForbidden("unable to delete non-own sharing token")
	}

	sharing.Active = false
	if err := SharingRepository.Save(&sharing); err != nil {
		return fmt.Errorf("fails to delete sharing: %w", err)
	}

	return nil
}

func GetSharing(token string) (Sharing, error) {
	sharing, err := SharingRepository.Find(token)
	if err != nil {
		return Sharing{}, fmt.Errorf("invalid sharing token: %w", err)
	}

	if !sharing.Active {
		return Sharing{}, util.NewErrorWithForbidden("sharing token has been inactive: %v", token)
	}

	return sharing, nil
}

/**
 * common
 */

func newSharingToken() string {
	return base64.RawStdEncoding.EncodeToString([]byte(uuid.NewString()))
}
