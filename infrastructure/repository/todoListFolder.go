package repository

import (
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type TodoListFolder struct {
	Entity

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

func (r TodoListFolderRepository) Save(entity *todolist.TodoListFolder) error {
	po := r.convertToPO(entity)
	re := r.db.Create(po).Error
	entity.ID = po.ID
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

func (r TodoListFolderRepository) Find(id int64) (todolist.TodoListFolder, error) {
	var folder TodoListFolder
	err := r.db.
		Where(&TodoListFolder{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&folder).
		Error

	return r.convertToEntity(folder), util.WrapGormErr(err, "todo list folder")
}

func (r TodoListFolderRepository) FindAllByUser(userId int64) ([]todolist.TodoListFolder, error) {
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

	entities := make([]todolist.TodoListFolder, len(POs))
	for i := range POs {
		entities = append(entities, r.convertToEntity(POs[i]))
	}

	return entities, nil
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

func (r TodoListFolderRepository) convertToPO(entity *todolist.TodoListFolder) TodoListFolder {
	return TodoListFolder{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},
		Name:      entity.Name,
		UserID:    entity.UserID,
		TodoLists: nil, // TODO
	}
}

func (r TodoListFolderRepository) convertToEntity(po TodoListFolder) todolist.TodoListFolder {
	return todolist.TodoListFolder{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
		Name:      po.Name,
		UserID:    po.UserID,
		TodoLists: nil, // TODO
	}
}
