package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
)

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

type ThirdPartyOAuthToken struct {
	Entity

	Active bool   `json:"active"`
	Type   int8   `json:"type" gorm:"index:idx_third_party_oauth_tokens_user,unique"`
	Token  string `json:"token" gorm:"size:128"`
	Scope  string `json:"scope" gorm:"size:32"`

	UserID int64 `json:"userID" gorm:"index:idx_third_party_oauth_tokens_user,unique"`
	User   User  `json:"-"`
}

func (ThirdPartyOAuthToken) TableName() string {
	return "third_party_oauth_tokens"
}

func InsertThirdPartyOAuthToken(entity *ThirdPartyOAuthToken) error {
	err := db.Create(entity).Error
	return util.WrapGormErr(err, "third party token")
}

func UpdateThirdPartyOAuthToken(new *ThirdPartyOAuthToken) error {
	err := db.
		Where(&ThirdPartyOAuthToken{
			UserID: new.UserID,
			Type:   new.Type,
		}).
		Save(new).
		Error

	return util.WrapGormErr(err, "third party token")
}

func ExistActiveThirdPartyOAuthToken(userID int64, tokenType ThirdPartyTokenType) (bool, error) {
	var count int64
	err := db.
		Model(ThirdPartyOAuthToken{}).
		Where(ThirdPartyOAuthToken{
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
