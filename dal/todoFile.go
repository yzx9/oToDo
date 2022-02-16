package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodoFile(file *entity.TodoFile) error {
	re := db.Create(file)
	return utils.WrapGormErr(re.Error, "todo file")
}

func SelectTodoFiles(todoID string) ([]entity.TodoFile, error) {
	var files []entity.TodoFile
	re := db.Where(entity.TodoFile{TodoID: todoID}).Find(&files)
	return files, utils.WrapGormErr(re.Error, "file")
}
