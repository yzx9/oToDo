package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoListFolder(todoListFolder *entity.TodoListFolder) error {
	re := db.Create(todoListFolder)
	return util.WrapGormErr(re.Error, "todo list folder")
}

func SelectTodoListFolder(id int64) (entity.TodoListFolder, error) {
	var folder entity.TodoListFolder
	where := entity.TodoListFolder{Entity: entity.Entity{ID: id}}
	re := db.Where(&where).First(&folder)
	return folder, util.WrapGormErr(re.Error, "todo list folder")
}

func SelectTodoListFolders(userId int64) ([]entity.TodoListFolder, error) {
	var folders []entity.TodoListFolder
	re := db.Where(entity.TodoListFolder{UserID: userId}).Find(&folders)
	return folders, util.WrapGormErr(re.Error, "todo list folder")
}

func DeleteTodoListFolder(id int64) error {
	re := db.Delete(&entity.TodoListFolder{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return util.WrapGormErr(re.Error, "todo list folder")
}

func ExistTodoListFolder(id int64) (bool, error) {
	var count int64
	folder := entity.TodoListFolder{Entity: entity.Entity{ID: id}}
	re := db.Model(&entity.TodoListFolder{}).Where(&folder).Count(&count)
	return count != 0, util.WrapGormErr(re.Error, "tag")
}
