package repository

import (
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type ThirdPartyOAuthToken struct {
	Entity

	Active bool
	Type   int8   `gorm:"index:idx_third_party_oauth_tokens_user,unique"`
	Token  string `gorm:"size:128"`
	Scope  string `gorm:"size:32"`

	UserID int64 `gorm:"index:idx_third_party_oauth_tokens_user,unique"`
}

func (ThirdPartyOAuthToken) TableName() string { return "third_party_oauth_tokens" }

type ThirdPartyOAuthTokenRepository struct {
	db *gorm.DB
}

func NewThirdPartyOAuthTokenRepository(db *gorm.DB) ThirdPartyOAuthTokenRepository {
	return ThirdPartyOAuthTokenRepository{db: db}
}

func (r ThirdPartyOAuthTokenRepository) Save(entity *identity.ThirdPartyOAuthToken) error {
	po := r.convertToPO(*entity)
	err := r.db.Save(&po).Error
	entity.SetID(po.ID)
	return util.WrapGormErr(err, "third party token")
}

// TODO: change primary key to userID+Type
func (r ThirdPartyOAuthTokenRepository) SaveByUserIDAndType(entity *identity.ThirdPartyOAuthToken) error {
	po := r.convertToPO(*entity)
	err := r.db.
		Where(map[string]any{
			"user_id": entity.UserID(),
			"type":    int8(entity.Type()),
		}).
		Save(&po).
		Error

	entity.SetID(po.ID)
	return util.WrapGormErr(err, "third party token")
}

func (r ThirdPartyOAuthTokenRepository) ExistActiveOne(userID int64, tokenType identity.ThirdPartyTokenType) (bool, error) {
	var count int64
	err := r.db.
		Model(ThirdPartyOAuthToken{}).
		Where(map[string]any{
			"user_id": userID,
			"type":    int8(tokenType),
			"active":  true,
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
			ID:        entity.ID(),
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},

		UserID: entity.UserID(),
		Active: entity.Active(),
		Type:   int8(entity.Type()),
		Token:  entity.Token(),
		Scope:  entity.Scope(),
	}
}

func (r ThirdPartyOAuthTokenRepository) convertToEntity(po ThirdPartyOAuthToken) identity.ThirdPartyOAuthToken {
	return identity.NewThirdPartyOAuthToken(
		po.ID,
		po.CreatedAt,
		po.UpdatedAt,
		po.UserID,
		po.Active,
		identity.ThirdPartyTokenType(po.Type),
		po.Token,
		po.Scope,
	)
}
