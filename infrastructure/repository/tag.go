package repository

import (
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type Tag struct {
	Entity

	Name string `gorm:"size:32;index:idx_tags_user,unique"`

	UserID int64 `gorm:"index:idx_tags_user,unique"`
	User   User

	Todos []Todo `gorm:"many2many:tag_todos;"`
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return TagRepository{db: db}
}

func (r TagRepository) Save(entity *todo.Tag) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "tag")
}

func (r TagRepository) Find(userID int64, tagName string) (todo.Tag, error) {
	var po Tag
	err := r.db.
		Scopes(filterTag(userID, tagName)).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "tag")
}

func (r TagRepository) FindAllByUser(userID int64) ([]todo.Tag, error) {
	var tags []Tag
	err := r.db.
		Where(Tag{
			UserID: userID,
		}).
		Find(&tags).
		Error

	return r.convertToEntities(tags), util.WrapGormErr(err, "tag")
}

func (r TagRepository) convertToPO(entity *todo.Tag) Tag {
	return Tag{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Name:   entity.Name,
		UserID: entity.UserID,
		Todos:  nil, // TODO
	}
}

func (r TagRepository) convertToEntity(po Tag) todo.Tag {
	return todo.Tag{
		ID:        po.ID,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,

		Name: po.Name,

		UserID: po.UserID,

		Todos: nil, // TODO
	}
}

func (r TagRepository) convertToEntities(POs []Tag) []todo.Tag {
	return util.Map(r.convertToEntity, POs)
}

type TagTodoRepository struct {
	db *gorm.DB
}

func NewTagTodoRepository(db *gorm.DB) TagTodoRepository {
	return TagTodoRepository{db: db}
}

func (r TagTodoRepository) Save(userID, todoID int64, tagName string) error {
	err := r.db.
		Scopes(filterTag(userID, tagName)).
		Association("Todos").
		Append(&Todo{
			Entity: Entity{
				ID: todoID,
			},
		})

	return util.WrapGormErr(err, "tag todos")
}

func (r TagTodoRepository) Delete(userID, todoID int64, tagName string) error {
	err := r.db.
		Scopes(filterTag(userID, tagName)).
		Association("Todos").
		Delete(&Todo{
			Entity: Entity{
				ID: todoID,
			},
		})

	return util.WrapGormErr(err, "tag todos")
}

func (r TagRepository) Exist(userID int64, tagName string) (bool, error) {
	var count int64
	err := r.db.
		Scopes(filterTag(userID, tagName)).
		Count(&count).
		Error

	if err != nil {
		return false, util.WrapGormErr(err, "tag")
	}

	return count != 0, nil
}

func filterTag(userID int64, tagName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Tag{}).
			Where(&Tag{
				UserID: userID,
				Name:   tagName,
			})
	}
}
