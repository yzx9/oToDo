package repository

import (
	"github.com/yzx9/otodo/domain/user"
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

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Save(entity *user.User) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.ID = po.ID
	return util.WrapGormErr(err, "user")
}

func (r UserRepository) Find(id int64) (user.User, error) {
	var po User
	err := r.db.
		Where(&User{
			Entity: Entity{
				ID: id,
			},
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByUserName(username string) (user.User, error) {
	var po User
	err := r.db.
		Where(User{
			Name: username,
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByGithubID(githubID int64) (user.User, error) {
	var po User
	err := r.db.
		Where(User{
			GithubID: githubID,
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByTodo(todoID int64) (user.User, error) {
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
		return user.User{}, util.WrapGormErr(err, "todo")
	}

	return r.Find(todo.UserID)
}

func (r UserRepository) ExistByUserName(username string) (bool, error) {
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

func (r UserRepository) ExistByGithubID(githubID int64) (bool, error) {
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

func (r UserRepository) convertToPO(entity *user.User) User {
	return User{
		Entity: Entity{
			ID:        entity.ID,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		},

		Name:      entity.Name,
		Nickname:  entity.Nickname,
		Password:  entity.Password,
		Email:     entity.Email,
		Telephone: entity.Telephone,
		Avatar:    entity.Avatar,
		GithubID:  entity.GithubID,

		BasicTodoListID: entity.BasicTodoListID,
		BasicTodoList:   nil, // TODO

		TodoLists: nil, // TODO

		SharedTodoLists: nil, // TODO
	}
}

func (r UserRepository) convertToEntity(po User) user.User {
	return user.User{
		ID:              po.ID,
		CreatedAt:       po.CreatedAt,
		UpdatedAt:       po.UpdatedAt,
		Name:            po.Name,
		Nickname:        po.Nickname,
		Password:        po.Password,
		Email:           po.Email,
		Telephone:       po.Telephone,
		Avatar:          po.Avatar,
		GithubID:        po.GithubID,
		BasicTodoListID: po.BasicTodoListID,
	}
}

func (r UserRepository) convertToEntities(POs []User) []user.User {
	return util.Map(r.convertToEntity, POs)
}
