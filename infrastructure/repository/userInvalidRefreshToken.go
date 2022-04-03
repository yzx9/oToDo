package repository

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertUserInvalidRefreshToken(entity *entity.UserInvalidRefreshToken) error {
	err := db.Create(entity).Error
	return util.WrapGormErr(err, "user invalid refresh token")
}

func ExistUserInvalidRefreshToken(userID int64, tokenID string) (bool, error) {
	var count int64
	err := db.
		Model(entity.UserInvalidRefreshToken{}).
		Where(&entity.UserInvalidRefreshToken{
			UserID:  userID,
			TokenID: tokenID,
		}).
		Count(&count).
		Error

	if err != nil {
		return false, util.WrapGormErr(err, "user invalid refresh token")
	}

	return count != 0, nil
}
