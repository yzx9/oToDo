package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
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

var SharingRepo SharingRepository

type SharingRepository struct {
	db *gorm.DB
}

func (r *SharingRepository) Save(sharing *Sharing) error {
	err := r.db.Save(&sharing).Error
	return util.WrapGormErr(err, "sharing")
}

func (r *SharingRepository) Find(token string) (Sharing, error) {
	var sharing Sharing
	err := r.db.
		Where(&Sharing{
			Token: token,
		}).
		First(&sharing).
		Error

	return sharing, util.WrapGormErr(err, "sharing")
}

func (r *SharingRepository) FindByUser(userID int64, sharingType SharingType) ([]Sharing, error) {
	var sharings []Sharing
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func (r *SharingRepository) FindAllActiveOnes(userID int64, sharingType SharingType) ([]Sharing, error) {
	var sharings []Sharing
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func (r *SharingRepository) ExistActiveOne(userID int64, sharingType SharingType) (bool, error) {
	var count int64
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "sharing")
}

func (r *SharingRepository) DeleteSharings(userID int64, sharingType SharingType) (int64, error) {
	// Here we inactive sharing instead of not delete
	re := r.db.
		Where(Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Update("active", false)

	return re.RowsAffected, util.WrapGormErr(re.Error, "sharing")
}
