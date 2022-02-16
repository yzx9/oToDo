package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertTag(tag *entity.Tag) error {
	re := db.Create(tag)
	return utils.WrapGormErr(re.Error, "tag")
}

func SelectTag(userID, tagName string) (entity.Tag, error) {
	var tag entity.Tag
	re := db.Where(&entity.Tag{
		UserID: userID,
		Name:   tagName,
	}).First(&tag)
	if re.Error != nil {
		return entity.Tag{}, utils.WrapGormErr(re.Error, "tag")
	}

	return tag, nil
}

func SelectTags(userID string) ([]entity.Tag, error) {
	var tags []entity.Tag
	re := db.Where(&entity.Tag{UserID: userID}).Find(&tags)
	if re.Error != nil {
		return nil, utils.WrapGormErr(re.Error, "tag")
	}

	return tags, nil
}

func InsertTagTodo(userID, todoID, tagName string) error {
	err := db.Model(&entity.Tag{}).Association("Todos").Append(&entity.Todo{
		Entity: entity.Entity{
			ID: todoID,
		},
	})

	return utils.WrapGormErr(err, "tag todos")
}

func DeleteTagTodo(userID, todoID, tagName string) error {
	err := db.Model(&entity.Tag{}).Association("Todos").Delete(&entity.Todo{
		Entity: entity.Entity{
			ID: todoID,
		},
	})

	return utils.WrapGormErr(err, "tag todos")
}

func ExistTag(userID, tagName string) (bool, error) {
	var count int64
	re := db.Model(&entity.Tag{}).Where(&entity.Tag{
		UserID: userID,
		Name:   tagName,
	}).Count(&count)
	if re.Error != nil {
		return false, utils.WrapGormErr(re.Error, "tag")
	}

	return count != 0, nil
}
