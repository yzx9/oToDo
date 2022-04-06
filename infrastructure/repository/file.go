package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type FileAccessType int

const (
	FileTypePublic FileAccessType = 10*iota + 1 // set RelatedID to empty
	FileTypeTodo                                // set RelatedID to TodoID
)

type File struct {
	Entity

	FileName     string `json:"fileName"`
	FileServerID string `json:"-" gorm:"size:15"`
	FilePath     string `json:"-" gorm:"size:128"`
	AccessType   int8   `json:"-"` // FileAccessType
	RelatedID    int64  `json:"-"` // Depend on access type
}

var FileRepo FileRepository

type FileRepository struct {
	db *gorm.DB
}

func (r *FileRepository) Insert(file *File) error {
	err := r.db.Create(file).Error
	return util.WrapGormErr(err, "file")
}

func (r *FileRepository) Save(file *File) error {
	err := r.db.Save(file).Error
	return util.WrapGormErr(err, "file")
}

func (r *FileRepository) Find(id int64) (*File, error) {
	var file File
	err := r.db.
		Where(&File{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&file).
		Error

	return &file, util.WrapGormErr(err, "file")
}
