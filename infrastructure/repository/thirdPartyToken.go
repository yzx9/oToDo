package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
)

func InsertThirdPartyOAuthToken(entity *entity.ThirdPartyOAuthToken) error {
	err := db.Create(entity).Error
	return util.WrapGormErr(err, "third party token")
}

func UpdateThirdPartyOAuthToken(new *entity.ThirdPartyOAuthToken) error {
	err := db.
		Where(&entity.ThirdPartyOAuthToken{
			UserID: new.UserID,
			Type:   new.Type,
		}).
		Save(new).
		Error

	return util.WrapGormErr(err, "third party token")
}

func ExistActiveThirdPartyOAuthToken(userID int64, tokenType entity.ThirdPartyTokenType) (bool, error) {
	var count int64
	err := db.
		Model(entity.ThirdPartyOAuthToken{}).
		Where(entity.ThirdPartyOAuthToken{
			UserID: userID,
			Type:   int8(tokenType),
			Active: true,
		}).
		Count(&count).
		Error

	if err != nil {
		return false, util.WrapGormErr(err, "third party token")
	}

	return count != 0, nil
}
