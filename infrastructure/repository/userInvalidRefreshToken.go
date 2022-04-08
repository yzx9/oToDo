package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type UserInvalidRefreshToken struct {
	Entity

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TokenID string `json:"tokenID" gorm:"type:char(36);"`
}

var UserInvalidRefreshTokenRepo UserInvalidRefreshTokenRepository

type UserInvalidRefreshTokenRepository struct {
	db *gorm.DB
}

func NewUserInvalidRefreshTokenRepository(db *gorm.DB) UserInvalidRefreshTokenRepository {
	return UserInvalidRefreshTokenRepository{db: db}
}

func (r UserInvalidRefreshTokenRepository) Save(entity *UserInvalidRefreshToken) error {
	err := r.db.Save(entity).Error
	return util.WrapGormErr(err, "user invalid refresh token")
}

func (r UserInvalidRefreshTokenRepository) Exist(userID int64, tokenID string) (bool, error) {
	var count int64
	err := r.db.
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
