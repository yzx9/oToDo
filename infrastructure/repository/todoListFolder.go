package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoListFolder struct {
	Entity

	Name string `json:"name" gorm:"size:128"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TodoLists []TodoList `json:"-"`
}

var TodoListFolderRepo TodoListFolderRepository

type TodoListFolderRepository struct {
	db *gorm.DB
}

func (r TodoListFolderRepository) Save(todoListFolder *TodoListFolder) error {
	re := r.db.Create(todoListFolder).Error
	return util.WrapGormErr(re, "todo list folder")
}

func (r TodoListFolderRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) Find(id int64) (TodoListFolder, error) {
	var folder TodoListFolder
	err := r.db.
		Where(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&folder).
		Error

	return folder, util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) FindAllByUser(userId int64) ([]TodoListFolder, error) {
	var folders []TodoListFolder
	err := r.db.
		Where(TodoListFolder{
			UserID: userId,
		}).
		Find(&folders).
		Error

	return folders, util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) Exist(id int64) (bool, error) {
	var count int64
	err := r.db.
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
