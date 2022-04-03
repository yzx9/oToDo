package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
)

type TodoListFolder struct {
	Entity

	Name string `json:"name" gorm:"size:128"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoLists []TodoList `json:"-"`
}

func InsertTodoListFolder(todoListFolder *TodoListFolder) error {
	re := db.Create(todoListFolder).Error
	return util.WrapGormErr(re, "todo list folder")
}

func SelectTodoListFolder(id int64) (TodoListFolder, error) {
	var folder TodoListFolder
	err := db.
		Where(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&folder).
		Error

	return folder, util.WrapGormErr(err, "todo list folder")
}

func SelectTodoListFolders(userId int64) ([]TodoListFolder, error) {
	var folders []TodoListFolder
	err := db.
		Where(TodoListFolder{
			UserID: userId,
		}).
		Find(&folders).
		Error

	return folders, util.WrapGormErr(err, "todo list folder")
}

func DeleteTodoListFolder(id int64) error {
	err := db.
		Delete(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list folder")
}

func ExistTodoListFolder(id int64) (bool, error) {
	var count int64
	err := db.
		Model(&TodoListFolder{}).
		Where(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "tag")
}
