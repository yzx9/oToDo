package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertFile(file *entity.File) error {
	re := db.Create(file)
	return util.WrapGormErr(re.Error, "file")
}

func SelectFile(id int64) (*entity.File, error) {
	var file entity.File
	where := entity.File{Entity: entity.Entity{ID: id}}
	re := db.Where(&where).First(&file)
	return &file, util.WrapGormErr(re.Error, "file")
}
