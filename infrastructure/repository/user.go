package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/model/entity"
)

func InsertUser(user *entity.User) error {
	re := db.Create(user)
	return util.WrapGormErr(re.Error, "user")
}

func SelectUser(id int64) (entity.User, error) {
	var user entity.User
	err := db.
		Where(&entity.User{
			Entity: entity.Entity{
				ID: id,
			},
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func SelectUserByUserName(username string) (entity.User, error) {
	var user entity.User
	err := db.
		Where(entity.User{
			Name: username,
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func SelectUserByGithubID(githubID int64) (entity.User, error) {
	var user entity.User
	err := db.
		Where(entity.User{
			GithubID: githubID,
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func SelectUserByTodo(todoID int64) (entity.User, error) {
	var todo entity.Todo
	err := db.
		Where(&entity.Todo{
			Entity: entity.Entity{
				ID: todoID,
			},
		}).
		Select("UserID").
		First(&todo).
		Error

	if err != nil {
		return entity.User{}, util.WrapGormErr(err, "todo")
	}

	return SelectUser(todo.UserID)
}

func SaveUser(user *entity.User) error {
	err := db.Save(&user).Error
	return util.WrapGormErr(err, "user")
}

func ExistUserByUserName(username string) (bool, error) {
	var count int64
	err := db.
		Model(&entity.User{}).
		Where(entity.User{
			Name: username,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "user")
}

func ExistUserByGithubID(githubID int64) (bool, error) {
	var count int64
	err := db.
		Model(&entity.User{}).
		Where(entity.User{
			GithubID: githubID,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "user")
}
