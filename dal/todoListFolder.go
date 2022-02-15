package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTodoListFolder(todoListFolder *entity.TodoListFolder) error {
	re := db.Create(todoListFolder)
	return utils.WrapGormErr(re.Error, "todo list folder")
}

func SelectTodoListFolder(id string) (entity.TodoListFolder, error) {
	var folder entity.TodoListFolder
	re := db.Where("ID = ?", id).First(&folder)
	return folder, utils.WrapGormErr(re.Error, "todo list folder")
}

func SelectTodoListFolders(userId string) ([]entity.TodoListFolder, error) {
	var folders []entity.TodoListFolder
	re := db.Where("UserID = ?", userId).Find(&folders)
	return folders, utils.WrapGormErr(re.Error, "todo list folder")
}

func DeleteTodoListFolder(id string) error {
	re := db.Delete(&entity.TodoListFolder{
		Entity: entity.Entity{
			ID: id,
		},
	})
	return utils.WrapGormErr(re.Error, "todo list folder")
}

func ExistTodoListFolder(id string) (bool, error) {
	var count int64
	re := db.Model(&entity.TodoListFolder{}).Where("ID = ?", id).Count(&count)
	return count != 0, utils.WrapGormErr(re.Error, "tag")
}
