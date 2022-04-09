package repository

import (
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type Sharing struct {
	Entity

	Token     string `gorm:"size:128;uniqueIndex"`
	Active    bool
	Type      int8  // SharingType
	RelatedID int64 // Depends on Type

	UserID int64
	User   User
}

var SharingRepo SharingRepository

type SharingRepository struct {
	db *gorm.DB
}

func NewSharingRepository(db *gorm.DB) SharingRepository {
	return SharingRepository{db: db}
}

func (r SharingRepository) Save(entity *todolist.Sharing) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "sharing")
}

func (r SharingRepository) DeleteAllByUserAndType(userID int64, sharingType todolist.SharingType) (int64, error) {
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

func (r SharingRepository) Find(token string) (todolist.Sharing, error) {
	var sharing Sharing
	err := r.db.
		Where(&Sharing{
			Token: token,
		}).
		First(&sharing).
		Error

	return r.convertToEntity(sharing), util.WrapGormErr(err, "sharing")
}

func (r SharingRepository) FindByUser(userID int64, sharingType todolist.SharingType) ([]todolist.Sharing, error) {
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

	entities := make([]todolist.Sharing, len(POs))
	for i := range POs {
		entities = append(entities, r.convertToEntity(POs[i]))
	}

	return entities, nil
}

func (r SharingRepository) FindAllActiveOnes(userID int64, sharingType todolist.SharingType) ([]todolist.Sharing, error) {
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

	entities := make([]todolist.Sharing, len(POs))
	for i := range POs {
		entities = append(entities, r.convertToEntity(POs[i]))
	}

	return entities, nil
}

func (r SharingRepository) ExistActiveOne(userID int64, sharingType todolist.SharingType) (bool, error) {
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

func (r SharingRepository) convertToPO(entity *todolist.Sharing) Sharing {
	return Sharing{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Token:     entity.Token,
		Active:    entity.Active,
		Type:      entity.Type,
		RelatedID: entity.RelatedID,

		UserID: entity.UserID,
	}
}

func (r SharingRepository) convertToEntity(po Sharing) todolist.Sharing {
	return todolist.Sharing{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Token:     po.Token,
		Active:    po.Active,
		Type:      po.Type,
		RelatedID: po.RelatedID,

		UserID: po.UserID,
	}
}
