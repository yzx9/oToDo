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

func InsertTag(tag *Tag) error {
	err := db.Create(tag).Error
	return util.WrapGormErr(err, "tag")
}

func SelectTag(userID int64, tagName string) (Tag, error) {
	var tag Tag
	err := db.
		Scopes(tagScope(userID, tagName)).
		First(&tag).
		Error

	return tag, util.WrapGormErr(err, "tag")
}

func SelectTags(userID int64) ([]Tag, error) {
	var tags []Tag
	err := db.
		Where(Tag{
			UserID: userID,
		}).
		Find(&tags).
		Error

	return tags, util.WrapGormErr(err, "tag")
}

func InsertTagTodo(userID, todoID int64, tagName string) error {
	err := db.
		Scopes(tagScope(userID, tagName)).
		Association("Todos").
		Append(&Todo{
			Entity: Entity{
				ID: todoID,
			},
		})

	return util.WrapGormErr(err, "tag todos")
}

func DeleteTagTodo(userID, todoID int64, tagName string) error {
	err := db.
		Scopes(tagScope(userID, tagName)).
		Association("Todos").
		Delete(&Todo{
			Entity: Entity{
				ID: todoID,
			},
		})

	return util.WrapGormErr(err, "tag todos")
}

func ExistTag(userID int64, tagName string) (bool, error) {
	var count int64
	err := db.
		Scopes(tagScope(userID, tagName)).
		Count(&count).
		Error

	if err != nil {
		return false, util.WrapGormErr(err, "tag")
	}

	return count != 0, nil
}

func tagScope(userID int64, tagName string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Model(&Tag{}).
			Where(&Tag{
				UserID: userID,
				Name:   tagName,
			})
	}
}
