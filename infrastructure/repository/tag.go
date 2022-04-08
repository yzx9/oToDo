package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type Tag struct {
	Entity

	Name string `json:"name" gorm:"size:32;index:idx_tags_user,unique"`

	UserID int64 `json:"userID" gorm:"index:idx_tags_user,unique"`
	User   User  `json:"-"`

	Todos []Todo `json:"-" gorm:"many2many:tag_todos;"`
}

var TagRepo TagRepository

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return TagRepository{db: db}
}

func (r TagRepository) Save(tag *Tag) error {
	err := r.db.Create(tag).Error
	return util.WrapGormErr(err, "tag")
}

func (r TagRepository) Find(userID int64, tagName string) (Tag, error) {
	var tag Tag
	err := r.db.
		Scopes(filterTag(userID, tagName)).
		First(&tag).
		Error

	return tag, util.WrapGormErr(err, "tag")
}

func (r TagRepository) FindAllByUser(userID int64) ([]Tag, error) {
	var tags []Tag
	err := r.db.
		Where(Tag{
			UserID: userID,
		}).
		Find(&tags).
		Error

	return tags, util.WrapGormErr(err, "tag")
}

var TagTodoRepo TagTodoRepository

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
