package repository

import (
	"github.com/yzx9/otodo/domain/sharing"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type Sharing struct {
	Entity

	Token     string `gorm:"size:128;uniqueIndex"`
	Active    bool
	Type      int8
	RelatedID int64

	UserID int64
	User   User
}

type SharingRepository struct {
	db *gorm.DB
}

func NewSharingRepository(db *gorm.DB) SharingRepository {
	return SharingRepository{db: db}
}

func (r SharingRepository) Save(entity *sharing.Sharing) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "sharing")
}

func (r SharingRepository) DeleteAllByUserAndType(userID int64, sharingType sharing.SharingType) (int64, error) {
	// Here we inactive sharing instead of not delete
	re := r.db.
		Where(Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Update("active", false)

	return re.RowsAffected, util.WrapGormErr(re.Error, "sharing")
}

func (r SharingRepository) Find(token string) (sharing.Sharing, error) {
	var sharing Sharing
	err := r.db.
		Where(&Sharing{
			Token: token,
		}).
		First(&sharing).
		Error

	return r.convertToEntity(sharing), util.WrapGormErr(err, "sharing")
}

func (r SharingRepository) FindByUser(userID int64, sharingType sharing.SharingType) ([]sharing.Sharing, error) {
	var POs []Sharing
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
		}).
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "sharing")
	}

	return r.convertToEntities(POs), nil
}

func (r SharingRepository) FindAllActive(userID int64, sharingType sharing.SharingType) ([]sharing.Sharing, error) {
	var POs []Sharing
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "sharing")
	}

	return r.convertToEntities(POs), nil
}

func (r SharingRepository) ExistActiveOne(userID int64, sharingType sharing.SharingType) (bool, error) {
	var count int64
	err := r.db.
		Where(&Sharing{
			UserID: userID,
			Type:   sharingType,
			Active: true,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "sharing")
}

func (r SharingRepository) convertToPO(entity *sharing.Sharing) Sharing {
	return Sharing{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Token:     entity.Token,
		Active:    entity.Active,
		Type:      int8(entity.Type),
		RelatedID: entity.RelatedID,

		UserID: entity.UserID,
	}
}

func (r SharingRepository) convertToEntity(po Sharing) sharing.Sharing {
	return sharing.Sharing{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Token:     po.Token,
		Active:    po.Active,
		Type:      sharing.SharingType(po.Type),
		RelatedID: po.RelatedID,

		UserID: po.UserID,
	}
}

func (r SharingRepository) convertToEntities(POs []Sharing) []sharing.Sharing {
	return util.Map(r.convertToEntity, POs)
}
