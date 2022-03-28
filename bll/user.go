package bll

import (
	"crypto/sha256"
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

func CreateUser(payload dto.CreateUserDTO) (entity.User, error) {
	if len(payload.UserName) < 5 {
		return entity.User{}, fmt.Errorf("user name too short: %v", payload.UserName)
	}

	if len(payload.Password) < 6 {
		return entity.User{}, fmt.Errorf("password too short")
	}

	exist, err := dal.ExistUserByUserName(payload.UserName)
	if err != nil {
		return entity.User{}, fmt.Errorf("fails to valid user name: %w", err)
	}

	if exist {
		return entity.User{}, util.NewError(otodo.ErrDuplicateID, "user name has been used: %v", payload.UserName)
	}

	user := entity.User{
		Name:     payload.UserName,
		Nickname: payload.Nickname,
		Password: GetCryptoPassword(payload.Password),
	}
	if err := createUser(&user); err != nil {
		return entity.User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

func GetUser(userID int64) (entity.User, error) {
	user, err := dal.SelectUser(userID)
	if err != nil {
		return entity.User{}, fmt.Errorf("fails to get user: %w", err)
	}

	return user, nil
}

/**
 * Invalid User Refresh Token
 */

func CreateUserInvalidRefreshToken(userID int64, tokenID string) (entity.UserInvalidRefreshToken, error) {
	model := entity.UserInvalidRefreshToken{
		UserID:  userID,
		TokenID: tokenID,
	}
	if err := dal.InsertUserInvalidRefreshToken(&model); err != nil {
		return entity.UserInvalidRefreshToken{}, fmt.Errorf("fails to make user refresh token invalid: %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID int64, tokenID string) (bool, error) {
	valid, err := dal.ExistUserInvalidRefreshToken(userID, tokenID)
	if err != nil {
		return false, fmt.Errorf("fails to get user refresh token: %w", err)
	}

	return valid, nil
}

// Password
func GetCryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), otodo.Conf.Secret.PasswordNonce...))
	return pwd[:]
}

/**
 * OAuth
 */

func getOrRegisterUserByGithub(profile dto.GithubUserPublicProfile) (entity.User, error) {
	exist, err := dal.ExistUserByGithubID(profile.ID)
	if err != nil {
		return entity.User{}, util.NewErrorWithUnknown("fails to register user: %w", err)
	}

	if exist {
		user, err := dal.SelectUserByGithubID(profile.ID)
		if err != nil {
			return entity.User{}, util.NewErrorWithUnknown("fails to get user: %w", err)
		}

		return user, nil
	}

	// Register new user
	// TODO[feat]: download user avatar
	user := entity.User{
		Name:     profile.Email,
		Nickname: profile.Name,
		Email:    profile.Email,
		GithubID: profile.ID,
	}
	if err := createUser(&user); err != nil {
		return entity.User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Helpers
 */

func createUser(user *entity.User) error {
	if err := dal.InsertUser(user); err != nil {
		return fmt.Errorf("fails to create user: %w", err)
	}

	// create base todo list
	if _, err := createBasicTodoList(user); err != nil {
		return fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return nil
}

func createBasicTodoList(user *entity.User) (entity.TodoList, error) {
	basicTodoList := entity.TodoList{
		Name:    "Todos", // TODO i18n
		IsBasic: true,
		UserID:  user.ID,
	}
	if err := dal.InsertTodoList(&basicTodoList); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	user.BasicTodoListID = basicTodoList.ID
	if err := dal.SaveUser(user); err != nil {
		return entity.TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return basicTodoList, nil
}
