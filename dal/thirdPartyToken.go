package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertThirdPartyOAuthToken(entity *entity.ThirdPartyOAuthToken) error {
	re := db.Create(entity)
	return util.WrapGormErr(re.Error, "third party token")
}

func UpdateThirdPartyOAuthToken(new *entity.ThirdPartyOAuthToken) error {
	re := db.
		Where(&entity.ThirdPartyOAuthToken{
			UserID: new.UserID,
			Type:   new.Type,
			Active: true,
		}).
		Save(new)

	return util.WrapGormErr(re.Error, "third party token")
}

func ExistActiveThirdPartyOAuthToken(userID int64, tokenType entity.ThirdPartyTokenType) (bool, error) {
	var count int64
	re := db.
		Where(&entity.ThirdPartyOAuthToken{
			UserID: userID,
			Type:   int8(tokenType),
			Active: true,
		}).
		Count(&count)
	if re.Error != nil {
		return false, util.WrapGormErr(re.Error, "third party token")
	}

	return count != 0, nil
}
