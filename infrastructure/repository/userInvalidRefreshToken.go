package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
)

type UserInvalidRefreshToken struct {
	Entity

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TokenID string `json:"tokenID" gorm:"type:char(36);"`
}

func InsertUserInvalidRefreshToken(entity *UserInvalidRefreshToken) error {
	err := db.Create(entity).Error
	return util.WrapGormErr(err, "user invalid refresh token")
}

func ExistUserInvalidRefreshToken(userID int64, tokenID string) (bool, error) {
	var count int64
	err := db.
		Model(UserInvalidRefreshToken{}).
		Where(&UserInvalidRefreshToken{
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
