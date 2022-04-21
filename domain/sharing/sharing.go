package sharing

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/util"
)

type Sharing struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Token     string
	Active    bool
	Type      SharingType
	RelatedID int64 // Depends on Type

	UserID int64
}

type SharingType = int8

const (
	SharingTypeTodoList SharingType = 10*iota + 1 // Set RelatedID to todo list id
)

func CreateSharing(userID, relatedID int64, sharingType SharingType) (Sharing, error) {
	// Only allow one sharing active
	if _, err := SharingRepository.DeleteAllByUserAndType(userID, SharingTypeTodoList); err != nil {
		return Sharing{}, fmt.Errorf("fails to delete old sharing tokens: %w", err)
	}

	sharing := Sharing{
		Token:     newSharingToken(),
		Active:    true,
		Type:      sharingType,
		RelatedID: relatedID,
		UserID:    userID,
	}
	if err := SharingRepository.Save(&sharing); err != nil {
		return Sharing{}, fmt.Errorf("fails to create sharing token: %w", err)
	}

	return sharing, nil
}

func (sharing Sharing) Delete() error {
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
