package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertSharing(sharing *entity.Sharing) error {
	re := db.Create(sharing)
	return util.WrapGormErr(re.Error, "sharing")
}

func SelectSharing(token string) (entity.Sharing, error) {
	var sharing entity.Sharing
	re := db.Where(&entity.Sharing{Token: token}).First(&sharing)
	return sharing, util.WrapGormErr(re.Error, "sharing")
}

func SelectSharings(userID int64, sharingType entity.SharingType) ([]entity.Sharing, error) {
	var sharings []entity.Sharing
	re := db.Where(&entity.Sharing{
		UserID: userID,
		Type:   sharingType,
	}).Find(&sharings)
	return sharings, util.WrapGormErr(re.Error, "sharing")
}

func SelectActiveSharings(userID int64, sharingType entity.SharingType) ([]entity.Sharing, error) {
	var sharings []entity.Sharing
	re := db.Where(&entity.Sharing{
		UserID: userID,
		Type:   sharingType,
		Active: true,
	}).Find(&sharings)
	return sharings, util.WrapGormErr(re.Error, "sharing")
}

func SaveSharing(sharing *entity.Sharing) error {
	re := db.Save(&sharing)
	return util.WrapGormErr(re.Error, "sharing")
}

func ExistActiveSharing(userID int64, sharingType entity.SharingType) (bool, error) {
	var count int64
	re := db.Where(&entity.Sharing{
		UserID: userID,
		Type:   sharingType,
		Active: true,
	}).Count(&count)
	return count != 0, util.WrapGormErr(re.Error, "sharing")
}

func DeleteSharings(userID int64, sharingType entity.SharingType) (int64, error) {
	// Here we inactive sharing instead of not delete
	re := db.Where(entity.Sharing{
		UserID: userID,
		Type:   sharingType,
		Active: true,
	}).Update("active", false)
	return re.RowsAffected, util.WrapGormErr(re.Error, "sharing")
}
