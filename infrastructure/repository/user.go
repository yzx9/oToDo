package repository

import (
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/gorm"
)

type User struct {
	Entity

	Name      string `json:"name" gorm:"size:128;index:,unique,priority:11;"`
	Nickname  string `json:"nickname" gorm:"size:128"`
	Password  []byte `json:"-" gorm:"size:32;"`
	Email     string `json:"email" gorm:"size:32;"`
	Telephone string `json:"telephone" gorm:"size:16;"`
	Avatar    string `json:"avatar"`
	GithubID  int64  `json:"githubID" gorm:"index:,unique,priority:12"`

	BasicTodoListID int64     `json:"basicTodoListID"`
	BasicTodoList   *TodoList `json:"-"`

	TodoLists []TodoList `json:"-"`

	SharedTodoLists []*TodoList `json:"-" gorm:"many2many:todo_list_shared_users"`
}

var UserRepo UserRepository

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) Insert(user *User) error {
	re := r.db.Create(user)
	return util.WrapGormErr(re.Error, "user")
}

func (r *UserRepository) Save(user *User) error {
	err := r.db.Save(&user).Error
	return util.WrapGormErr(err, "user")
}

func (r *UserRepository) SelectUser(id int64) (User, error) {
	var user User
	err := r.db.
		Where(&User{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func (r *UserRepository) FindByUserName(username string) (User, error) {
	var user User
	err := r.db.
		Where(User{
			Name: username,
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func (r *UserRepository) FindByGithubID(githubID int64) (User, error) {
	var user User
	err := r.db.
		Where(User{
			GithubID: githubID,
		}).
		First(&user).
		Error

	return user, util.WrapGormErr(err, "user")
}

func (r *UserRepository) FindByTodo(todoID int64) (User, error) {
	var todo Todo
	err := r.db.
		Where(&Todo{
			Entity: Entity{
				ID: todoID,
			},
		}).
		Select("UserID").
		First(&todo).
		Error

	if err != nil {
		return User{}, util.WrapGormErr(err, "todo")
	}

	return r.SelectUser(todo.UserID)
}

func (r *UserRepository) ExistByUserName(username string) (bool, error) {
	var count int64
	err := r.db.
		Model(&User{}).
		Where(User{
			Name: username,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "user")
}

func (r *UserRepository) ExistByGithubID(githubID int64) (bool, error) {
	var count int64
	err := r.db.
		Model(&User{}).
		Where(User{
			GithubID: githubID,
		}).
		Count(&count).
		Error

	return count != 0, util.WrapGormErr(err, "user")
}
