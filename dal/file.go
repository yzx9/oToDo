package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertFile(file *entity.File) error {
	re := db.Create(file)
	if re.Error != nil {
		return utils.WrapGormErr(re.Error, "file")
	}

	return nil
}

func SelectFile(id string) (entity.File, error) {
	var file entity.File
	re := db.Where("ID = ?", id).First(&file)
	if re.Error != nil {
		return entity.File{}, utils.WrapGormErr(re.Error, "file")
	}

	return file, nil
}
