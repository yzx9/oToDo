package identity

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/util"
)

type User struct {
	ID              int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Name            string
	Nickname        string
	Password        []byte
	Email           string
	Telephone       string
	Avatar          string
	GithubID        int64
	BasicTodoListID int64
}

type NewUser struct {
	UserName string
	Password string
	Nickname string
}

func CreateUser(payload NewUser) (User, error) {
	if len(payload.UserName) < 5 {
		return User{}, fmt.Errorf("user name too short: %v", payload.UserName)
	}

	if len(payload.Password) < 6 {
		return User{}, fmt.Errorf("password too short")
	}

	exist, err := UserRepository.ExistByUserName(payload.UserName)
	if err != nil {
		return User{}, fmt.Errorf("fails to valid user name: %w", err)
	}

	if exist {
		return User{}, util.NewError(errors.ErrDuplicateID, "user name has been used: %v", payload.UserName)
	}

	user := User{
		Name:     payload.UserName,
		Nickname: payload.Nickname,
		Password: GetCryptoPassword(payload.Password),
	}
	if err := user.new(); err != nil {
		return User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

// Password
func (user User) IsSamePassword(password string) bool {
	crypto := GetCryptoPassword(password)
	return bytes.Equal(user.Password, crypto)
}

func GetCryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), config.Secret.PasswordNonce...))
	return pwd[:]
}

/**
 * OAuth
 */

func GetOrRegisterUserByGithub(profile GithubUserPublicProfile) (User, error) {
	exist, err := UserRepository.ExistByGithubID(profile.ID)
	if err != nil {
		return User{}, util.NewErrorWithUnknown("fails to register user: %w", err)
	}

	if exist {
		user, err := UserRepository.FindByGithubID(profile.ID)
		if err != nil {
			return User{}, util.NewErrorWithUnknown("fails to get user: %w", err)
		}

		return user, nil
	}

	// Register new user
	// TODO[feat]: download user avatar
	user := User{
		Name:     profile.Email,
		Nickname: profile.Name,
		Email:    profile.Email,
		GithubID: profile.ID,
	}
	if err := user.new(); err != nil {
		return User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Helpers
 */

func (user *User) new() error {
	if err := UserRepository.Save(user); err != nil {
		return fmt.Errorf("fails to create user: %w", err)
	}

	// create base todo list
	if _, err := user.createBasicTodoList(); err != nil {
		return fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return nil
}

// TODO move to todo list domain
func (user *User) createBasicTodoList() (todolist.TodoList, error) {
	basicTodoList := todolist.TodoList{
		Name:    "Todos", // TODO i18n
		IsBasic: true,
		UserID:  user.ID,
	}
	if err := TodoListRepo.Save(&basicTodoList); err != nil {
		return todolist.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	user.BasicTodoListID = basicTodoList.ID
	if err := UserRepository.Save(user); err != nil {
		return todolist.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return basicTodoList, nil
}
