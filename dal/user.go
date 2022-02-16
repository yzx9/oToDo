package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/utils"
)

func InsertUser(user *entity.User) error {
	re := db.Create(user)
	return utils.WrapGormErr(re.Error, "user")
}

func SelectUser(id string) (entity.User, error) {
	var user entity.User
	re := db.Where("id = ?", id).First(&user)
	return user, utils.WrapGormErr(re.Error, "user")
}

func SelectUserByUserName(username string) (entity.User, error) {
	var user entity.User
	re := db.Where(entity.User{Name: username}).First(&user)
	return user, utils.WrapGormErr(re.Error, "user")
}

func SelectUserByTodo(todoID string) (entity.User, error) {
	var todo entity.Todo
	re := db.Where("id = ?", todoID).Select("UserID").First(&todo)
	if re.Error != nil {
		return entity.User{}, utils.WrapGormErr(re.Error, "todo")
	}

	return SelectUser(todo.UserID)
}

func SaveUser(user *entity.User) error {
	re := db.Save(&user)
	return utils.WrapGormErr(re.Error, "user")
}

func ExistUserByUserName(username string) (bool, error) {
	var count int64
	re := db.Model(&entity.User{}).Where(entity.User{Name: username}).Count(&count)
	return count != 0, utils.WrapGormErr(re.Error, "user")
}
