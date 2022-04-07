package user

import (
	"crypto/sha256"
	"fmt"

	"github.com/yzx9/otodo/acl/github"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

type User struct {
	repository.Entity

	Name      string
	Nickname  string
	Password  []byte
	Email     string
	Telephone string
	Avatar    string
	GithubID  int64

	BasicTodoListID int64
	BasicTodoList   *repository.TodoList

	TodoLists []repository.TodoList

	SharedTodoLists []*repository.TodoList
}

type NewUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func CreateUser(payload NewUser) (repository.User, error) {
	if len(payload.UserName) < 5 {
		return repository.User{}, fmt.Errorf("user name too short: %v", payload.UserName)
	}

	if len(payload.Password) < 6 {
		return repository.User{}, fmt.Errorf("password too short")
	}

	exist, err := repository.UserRepo.ExistByUserName(payload.UserName)
	if err != nil {
		return repository.User{}, fmt.Errorf("fails to valid user name: %w", err)
	}

	if exist {
		return repository.User{}, util.NewError(errors.ErrDuplicateID, "user name has been used: %v", payload.UserName)
	}

	user := repository.User{
		Name:     payload.UserName,
		Nickname: payload.Nickname,
		Password: GetCryptoPassword(payload.Password),
	}
	if err := createUser(&user); err != nil {
		return repository.User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Invalid User Refresh Token
 */

func CreateUserInvalidRefreshToken(userID int64, tokenID string) (repository.UserInvalidRefreshToken, error) {
	model := repository.UserInvalidRefreshToken{
		UserID:  userID,
		TokenID: tokenID,
	}
	if err := repository.UserInvalidRefreshTokenRepo.Save(&model); err != nil {
		return repository.UserInvalidRefreshToken{}, fmt.Errorf("fails to make user refresh token invalid: %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID int64, tokenID string) (bool, error) {
	valid, err := repository.UserInvalidRefreshTokenRepo.Exist(userID, tokenID)
	if err != nil {
		return false, fmt.Errorf("fails to get user refresh token: %w", err)
	}

	return valid, nil
}

// Password
func GetCryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), config.Secret.PasswordNonce...))
	return pwd[:]
}

/**
 * OAuth
 */

func getOrRegisterUserByGithub(profile github.UserPublicProfile) (repository.User, error) {
	exist, err := repository.UserRepo.ExistByGithubID(profile.ID)
	if err != nil {
		return repository.User{}, util.NewErrorWithUnknown("fails to register user: %w", err)
	}

	if exist {
		user, err := repository.UserRepo.FindByGithubID(profile.ID)
		if err != nil {
			return repository.User{}, util.NewErrorWithUnknown("fails to get user: %w", err)
		}

		return user, nil
	}

	// Register new user
	// TODO[feat]: download user avatar
	user := repository.User{
		Name:     profile.Email,
		Nickname: profile.Name,
		Email:    profile.Email,
		GithubID: profile.ID,
	}
	if err := createUser(&user); err != nil {
		return repository.User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Helpers
 */

func createUser(user *repository.User) error {
	if err := repository.UserRepo.Save(user); err != nil {
		return fmt.Errorf("fails to create user: %w", err)
	}

	// create base todo list
	if _, err := createBasicTodoList(user); err != nil {
		return fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return nil
}

func createBasicTodoList(user *repository.User) (repository.TodoList, error) {
	basicTodoList := repository.TodoList{
		Name:    "Todos", // TODO i18n
		IsBasic: true,
		UserID:  user.ID,
	}
	if err := repository.TodoListRepo.Insert(&basicTodoList); err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	user.BasicTodoListID = basicTodoList.ID
	if err := repository.UserRepo.Save(user); err != nil {
		return repository.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return basicTodoList, nil
}
