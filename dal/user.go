package dal

import (
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertUser(user *entity.User) error {
	re := db.Create(user)
	return util.WrapGormErr(re.Error, "user")
}

func SelectUser(id int64) (entity.User, error) {
	var user entity.User
	where := entity.User{Entity: entity.Entity{ID: id}}
	re := db.Where(&where).First(&user)
	return user, util.WrapGormErr(re.Error, "user")
}

func SelectUserByUserName(username string) (entity.User, error) {
	var user entity.User
	re := db.Where(entity.User{Name: username}).First(&user)
	return user, util.WrapGormErr(re.Error, "user")
}

func SelectUserByTodo(todoID int64) (entity.User, error) {
	var todo entity.Todo
	where := entity.Todo{Entity: entity.Entity{ID: todoID}}
	re := db.Where(&where).Select("UserID").First(&todo)
	if re.Error != nil {
		return entity.User{}, util.WrapGormErr(re.Error, "todo")
	}

	return SelectUser(todo.UserID)
}

func SaveUser(user *entity.User) error {
	re := db.Save(&user)
	return util.WrapGormErr(re.Error, "user")
}

func ExistUserByUserName(username string) (bool, error) {
	var count int64
	re := db.Model(&entity.User{}).Where(entity.User{Name: username}).Count(&count)
	return count != 0, util.WrapGormErr(re.Error, "user")
}
