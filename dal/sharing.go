package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertSharing(sharing *entity.Sharing) error {
	err := db.Create(sharing).Error
	return util.WrapGormErr(err, "sharing")
}

func SelectSharing(token string) (entity.Sharing, error) {
	var sharing entity.Sharing
	err := db.
		Where(&entity.Sharing{
			Token: token,
		}).
		First(&sharing).
		Error

	return sharing, util.WrapGormErr(err, "sharing")
}

func SelectSharings(userID int64, sharingType entity.SharingType) ([]entity.Sharing, error) {
	var sharings []entity.Sharing
	err := db.
		Where(&entity.Sharing{
			UserID: userID,
			Type:   sharingType,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func SelectActiveSharings(userID int64, sharingType entity.SharingType) ([]entity.Sharing, error) {
	var sharings []entity.Sharing
	err := db.
		Where(&entity.Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Find(&sharings).
		Error

	return sharings, util.WrapGormErr(err, "sharing")
}

func SaveSharing(sharing *entity.Sharing) error {
	err := db.Save(&sharing).Error
	return util.WrapGormErr(err, "sharing")
}

func ExistActiveSharing(userID int64, sharingType entity.SharingType) (bool, error) {
	var count int64
	err := db.
		Where(&entity.Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "sharing")
}

func DeleteSharings(userID int64, sharingType entity.SharingType) (int64, error) {
	// Here we inactive sharing instead of not delete
	re := db.
		Where(entity.Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Update("active", false)

	return re.RowsAffected, util.WrapGormErr(re.Error, "sharing")
}
