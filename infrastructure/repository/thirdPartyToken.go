package repository

import (
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
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

type ThirdPartyOAuthTokenRepository struct {
	db *gorm.DB
}

func NewThirdPartyOAuthTokenRepository(db *gorm.DB) ThirdPartyOAuthTokenRepository {
	return ThirdPartyOAuthTokenRepository{db: db}
}

func (r ThirdPartyOAuthTokenRepository) Save(entity *identity.ThirdPartyOAuthToken) error {
	po := r.convertToPO(*entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "third party token")
}

// TODO: change primary key to userID+Type
func (r ThirdPartyOAuthTokenRepository) SaveByUserIDAndType(entity *identity.ThirdPartyOAuthToken) error {
	po := r.convertToPO(*entity)
	err := r.db.
		Where(&ThirdPartyOAuthToken{
			UserID: entity.UserID,
			Type:   int8(entity.Type),
		}).
		Save(&po).
		Error

	entity.ID = po.ID
	return util.WrapGormErr(err, "third party token")
}

func (r ThirdPartyOAuthTokenRepository) ExistActiveOne(userID int64, tokenType identity.ThirdPartyTokenType) (bool, error) {
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

func (r ThirdPartyOAuthTokenRepository) convertToPO(entity identity.ThirdPartyOAuthToken) ThirdPartyOAuthToken {
	return ThirdPartyOAuthToken{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Active: entity.Active,
		Type:   int8(entity.Type),
		Token:  entity.Token,
		Scope:  entity.Scope,

		UserID: entity.UserID,
	}
}

func (r ThirdPartyOAuthTokenRepository) convertToEntity(po ThirdPartyOAuthToken) identity.ThirdPartyOAuthToken {
	return identity.ThirdPartyOAuthToken{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Active: po.Active,
		Type:   identity.ThirdPartyTokenType(po.Type),
		Token:  po.Token,
		Scope:  po.Scope,

		UserID: po.UserID,
	}
}
