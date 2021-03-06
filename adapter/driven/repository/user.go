package repository

import (
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/util"
	"gorm.io/gorm"
)

type User struct {
	repository.Entity

	Name      string `json:"name" gorm:"size:128;index:,unique,priority:11;"`
	Nickname  string `json:"nickname" gorm:"size:128"`
	Password  []byte `json:"-" gorm:"size:32;"`
	Email     string `json:"email" gorm:"size:32;"`
	Telephone string `json:"telephone" gorm:"size:16;"`
	Avatar    string `json:"avatar"`
	GithubID  int64  `json:"githubID" gorm:"index:,unique,priority:12"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Save(entity *identity.User) error {
	po := r.convertToPO(entity)
	err := r.db.Save(&po).Error
	entity.SetID(po.ID)
	return util.WrapGormErr(err, "user")
}

func (r UserRepository) Find(id int64) (identity.User, error) {
	var po User
	err := r.db.
		Where(&User{
			Entity: repository.Entity{
				ID: id,
			},
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByUserName(username string) (identity.User, error) {
	var po User
	err := r.db.
		Where(User{
			Name: username,
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByGithubID(githubID int64) (identity.User, error) {
	var po User
	err := r.db.
		Where(User{
			GithubID: githubID,
		}).
		First(&po).
		Error

	return r.convertToEntity(po), util.WrapGormErr(err, "user")
}

func (r UserRepository) FindByTodo(todoID int64) (identity.User, error) {
	var todo Todo
	err := r.db.
		Where(&Todo{
			Entity: repository.Entity{
				ID: todoID,
			},
		}).
		Select("UserID").
		First(&todo).
		Error

	if err != nil {
		return identity.User{}, util.WrapGormErr(err, "todo")
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

func (r UserRepository) convertToPO(entity *identity.User) User {
	return User{
		Entity: repository.Entity{
			ID:        entity.Id(),
			CreatedAt: entity.CreatedAt(),
			UpdatedAt: entity.UpdatedAt(),
		},

		Name:      entity.Name(),
		Nickname:  entity.Nickname(),
		Password:  entity.Password().Bytes(),
		Email:     entity.Email(),
		Telephone: entity.Telephone(),
		Avatar:    entity.Avatar(),
		GithubID:  entity.GithubId(),
	}
}

func (r UserRepository) convertToEntity(po User) identity.User {
	return identity.NewUser(
		po.ID,
		po.CreatedAt,
		po.UpdatedAt,
		po.Name,
		po.Nickname,
		identity.NewPasswordByBytes(po.Password),
		po.Email,
		po.Telephone,
		po.Avatar,
		po.GithubID,
	)
}

func (r UserRepository) convertToEntities(POs []User) []identity.User {
	return util.Map(r.convertToEntity, POs)
}
