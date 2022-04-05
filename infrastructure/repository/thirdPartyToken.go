package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
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

var ThirdPartyOAuthTokenRepo ThirdPartyOAuthTokenRepository

type ThirdPartyOAuthTokenRepository struct {
	db *gorm.DB
}

func (r *ThirdPartyOAuthTokenRepository) InsertThirdPartyOAuthToken(entity *ThirdPartyOAuthToken) error {
	err := r.db.Create(entity).Error
	return util.WrapGormErr(err, "third party token")
}

func (r *ThirdPartyOAuthTokenRepository) UpdateThirdPartyOAuthToken(new *ThirdPartyOAuthToken) error {
	err := r.db.
		Where(&ThirdPartyOAuthToken{
			UserID: new.UserID,
			Type:   new.Type,
		}).
		Save(new).
		Error

	return util.WrapGormErr(err, "third party token")
}

func (r *ThirdPartyOAuthTokenRepository) ExistActiveThirdPartyOAuthToken(userID int64, tokenType ThirdPartyTokenType) (bool, error) {
	var count int64
	err := r.db.
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
