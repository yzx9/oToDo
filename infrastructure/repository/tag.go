package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
	"gorm.io/gorm"
)

func InsertTag(tag *entity.Tag) error {
	err := db.Create(tag).Error
	return util.WrapGormErr(err, "tag")
}

func SelectTag(userID int64, tagName string) (entity.Tag, error) {
	var tag entity.Tag
	err := db.
		Scopes(tagScope(userID, tagName)).
		First(&tag).
		Error

	return tag, util.WrapGormErr(err, "tag")
}

func SelectTags(userID int64) ([]entity.Tag, error) {
	var tags []entity.Tag
	err := db.
		Where(entity.Tag{
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
		Append(&entity.Todo{
			Entity: entity.Entity{
				ID: todoID,
			},
		})

	return util.WrapGormErr(err, "tag todos")
}

func DeleteTagTodo(userID, todoID int64, tagName string) error {
	err := db.
		Scopes(tagScope(userID, tagName)).
		Association("Todos").
		Delete(&entity.Todo{
			Entity: entity.Entity{
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
			Model(&entity.Tag{}).
			Where(&entity.Tag{
				UserID: userID,
				Name:   tagName,
			})
	}
}
