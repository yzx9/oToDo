package repository

import (
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type File struct {
	Entity

	FileName     string
	FileServerID string `gorm:"size:15"`
	FilePath     string `gorm:"size:128"`
	AccessType   int8   // FileAccessType
	RelatedID    int64  // Depend on access type
}

type FileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return FileRepository{db: db}
}

func (r FileRepository) Save(f *file.File) error {
	po := r.convertToPO(f)
	err := r.db.Save(&po).Error
	f.ID = po.ID
	return util.WrapGormErr(err, "file")
}

func (r FileRepository) Find(id int64) (*file.File, error) {
	var po File
	err := r.db.
		Where(&File{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&po).
		Error

	entity := r.convertToEntity(po)
	return &entity, util.WrapGormErr(err, "file")
}

// TODO mapper
func (r FileRepository) convertToPO(f *file.File) File {
	return File{
		Entity: Entity{
			ID:        f.ID,
			CreatedAt: f.CreatedAt,
			UpdatedAt: f.UpdatedAt,
		},
		FileName:     f.FileName,
		FileServerID: f.FileServerID,
		FilePath:     f.FilePath,
		AccessType:   int8(f.AccessType),
		RelatedID:    f.RelatedID,
	}
}

// TODO mapper
func (r FileRepository) convertToEntity(f File) file.File {
	return file.File{
		ID:           f.ID,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
		FileName:     f.FileName,
		FileServerID: f.FileServerID,
		FilePath:     f.FilePath,
		AccessType:   file.FileAccessType(f.AccessType),
		RelatedID:    f.RelatedID,
	}
}
