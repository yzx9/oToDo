package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertUserInvalidRefreshToken(entity entity.UserInvalidRefreshToken) error {
	re := db.Create(entity)
	if re.Error != nil {
		return utils.WrapGormErr(re.Error, "user invalid refresh token")
	}

	return nil
}

func ExistUserInvalidRefreshToken(userID, tokenID string) (bool, error) {
	var count int64
	re := db.Where(&entity.UserInvalidRefreshToken{
		Entity: entity.Entity{
			Deleted: false,
		},
		UserID:  userID,
		TokenID: tokenID,
	}).Count(&count)
	if re.Error != nil {
		return false, utils.WrapGormErr(re.Error, "user invalid refresh token")
	}

	return count != 0, nil
}
