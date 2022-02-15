package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertFile(file *entity.File) error {
	re := db.Create(file)
	return utils.WrapGormErr(re.Error, "file")
}

func SelectFile(id string) (*entity.File, error) {
	var file entity.File
	re := db.Where("id = ?", id).First(&file)
	return &file, utils.WrapGormErr(re.Error, "file")
}
