package repository

import (
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type TodoListFolder struct {
	repository.Entity

	Name string `gorm:"size:128"`

	UserID int64
	User   User

	TodoLists []TodoList
}

type TodoListFolderRepository struct {
	db *gorm.DB
}

func NewTodoListFolderRepository(db *gorm.DB) TodoListFolderRepository {
	return TodoListFolderRepository{db: db}
}

func (r TodoListFolderRepository) Save(entity *todo.TodoListFolder) error {
	po := r.convertToPO(entity)
	re := r.db.Create(po).Error
	entity.ID = po.ID
	return util.WrapGormErr(re, "todo list folder")
}

func (r TodoListFolderRepository) Delete(id int64) error {
	err := r.db.
		Delete(&TodoListFolder{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		Error

	return util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) Find(id int64) (todo.TodoListFolder, error) {
	var folder TodoListFolder
	err := r.db.
		Where(&TodoListFolder{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		First(&folder).
		Error

	return r.convertToEntity(folder), util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) FindAllByUser(userId int64) ([]todo.TodoListFolder, error) {
	var POs []TodoListFolder
	err := r.db.
		Where(TodoListFolder{
			UserID: userId,
		}).
		Find(&POs).
		Error

	if err != nil {
		return nil, util.WrapGormErr(err, "todo list folder")
	}

	return r.convertToEntities(POs), nil
}

func (r TodoListFolderRepository) Exist(id int64) (bool, error) {
	var count int64
	err := r.db.
		Model(&TodoListFolder{}).
		Where(&TodoListFolder{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "tag")
}

func (r TodoListFolderRepository) convertToPO(entity *todo.TodoListFolder) TodoListFolder {
	return TodoListFolder{
		Entity: repository.Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:      entity.Name,
		UserID:    entity.UserID,
		TodoLists: nil, // TODO
	}
}

func (r TodoListFolderRepository) convertToEntity(po TodoListFolder) todo.TodoListFolder {
	return todo.TodoListFolder{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		Name:      po.Name,
		UserID:    po.UserID,
		TodoLists: nil, // TODO
	}
}

func (r TodoListFolderRepository) convertToEntities(POs []TodoListFolder) []todo.TodoListFolder {
	return util.Map(r.convertToEntity, POs)
}
