package repository

import (
	"github.com/yzx9/otodo/domain/session"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type UserInvalidRefreshToken struct {
	Entity

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TokenID string `json:"tokenID" gorm:"type:char(36);"`
}

type UserInvalidRefreshTokenRepository struct {
	db *gorm.DB
}

func NewUserInvalidRefreshTokenRepository(db *gorm.DB) UserInvalidRefreshTokenRepository {
	return UserInvalidRefreshTokenRepository{db: db}
}

func (r UserInvalidRefreshTokenRepository) Save(entity *session.UserInvalidRefreshToken) error {
	po := r.convertToPO(*entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
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

func (r UserInvalidRefreshTokenRepository) convertToPO(entity session.UserInvalidRefreshToken) UserInvalidRefreshToken {
	return UserInvalidRefreshToken{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		UserID:  entity.UserID,
		TokenID: entity.TokenID,
	}
}

func (r UserInvalidRefreshTokenRepository) convertToEntity(po UserInvalidRefreshToken) session.UserInvalidRefreshToken {
	return session.UserInvalidRefreshToken{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		UserID:    po.UserID,
		TokenID:   po.TokenID,
	}
}
