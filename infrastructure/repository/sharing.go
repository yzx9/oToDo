package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
)

type SharingType = int8

const (
	SharingTypeTodoList SharingType = 10*iota + 1 // Set RelatedID to todo list id
)

type Sharing struct {
	Entity

	Token     string `json:"-" gorm:"size:128;uniqueIndex"`
	Active    bool   `json:"active"`
	Type      int8   `json:"type"` // SharingType
	RelatedID int64  `json:"-"`    // Depends on Type

	UserID int64 `json:"-"`
	User   User  `json:"-"`
}

func InsertSharing(sharing *Sharing) error {
	err := db.Create(sharing).Error
	return util.WrapGormErr(err, "sharing")
}

func SelectSharing(token string) (Sharing, error) {
	var sharing Sharing
	err := db.
		Where(&Sharing{
			Token: token,
		}).
		First(&sharing).
		Error

	return sharing, util.WrapGormErr(err, "sharing")
}

func SelectSharings(userID int64, sharingType SharingType) ([]Sharing, error) {
	var sharings []Sharing
	err := db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func SelectActiveSharings(userID int64, sharingType SharingType) ([]Sharing, error) {
	var sharings []Sharing
	err := db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func SaveSharing(sharing *Sharing) error {
	err := db.Save(&sharing).Error
	return util.WrapGormErr(err, "sharing")
}

func ExistActiveSharing(userID int64, sharingType SharingType) (bool, error) {
	var count int64
	err := db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "sharing")
}

func DeleteSharings(userID int64, sharingType SharingType) (int64, error) {
	// Here we inactive sharing instead of not delete
	re := db.
		Where(Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Update("active", false)

	return re.RowsAffected, util.WrapGormErr(re.Error, "sharing")
}
