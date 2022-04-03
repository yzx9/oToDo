package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
)

func InsertFile(file *entity.File) error {
	err := db.Create(file).Error
	return util.WrapGormErr(err, "file")
}

func SelectFile(id int64) (*entity.File, error) {
	var file entity.File
	err := db.
		Where(&entity.File{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&file).
		Error

	return &file, util.WrapGormErr(err, "file")
}

func SaveFile(file *entity.File) error {
	err := db.Save(file).Error
	return util.WrapGormErr(err, "file")
}
