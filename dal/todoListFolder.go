package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoListFolder(todoListFolder *entity.TodoListFolder) error {
	re := db.Create(todoListFolder).Error
	return util.WrapGormErr(re, "todo list folder")
}

func SelectTodoListFolder(id int64) (entity.TodoListFolder, error) {
	var folder entity.TodoListFolder
	err := db.
		Where(&entity.TodoListFolder{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&folder).
		Error

	return folder, util.WrapGormErr(err, "todo list folder")
}

func SelectTodoListFolders(userId int64) ([]entity.TodoListFolder, error) {
	var folders []entity.TodoListFolder
	err := db.
		Where(entity.TodoListFolder{
			UserID: userId,
		}).
		Find(&folders).
		Error

	return folders, util.WrapGormErr(err, "todo list folder")
}

func DeleteTodoListFolder(id int64) error {
	err := db.
		Delete(&entity.TodoListFolder{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list folder")
}

func ExistTodoListFolder(id int64) (bool, error) {
	var count int64
	err := db.
		Model(&entity.TodoListFolder{}).
		Where(&entity.TodoListFolder{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "tag")
}
