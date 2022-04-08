package user

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/yzx9/otodo/acl/github"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Name      string
	Nickname  string
	Password  []byte
	Email     string
	Telephone string
	Avatar    string
	GithubID  int64

	BasicTodoListID int64
	BasicTodoList   *TodoList

	TodoLists []TodoList

	SharedTodoLists []*TodoList
}

type NewUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
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
	if err := createUser(&user); err != nil {
		return User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Invalid User Refresh Token
 */
type UserInvalidRefreshToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID  int64
	TokenID string
}

func CreateUserInvalidRefreshToken(userID int64, tokenID string) (UserInvalidRefreshToken, error) {
	model := UserInvalidRefreshToken{
		UserID:  userID,
		TokenID: tokenID,
	}
	if err := UserInvalidRefreshTokenRepository.Save(&model); err != nil {
		return UserInvalidRefreshToken{}, fmt.Errorf("fails to make user refresh token invalid: %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID int64, tokenID string) (bool, error) {
	valid, err := UserInvalidRefreshTokenRepository.Exist(userID, tokenID)
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

func getOrRegisterUserByGithub(profile github.UserPublicProfile) (User, error) {
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
	if err := createUser(&user); err != nil {
		return User{}, fmt.Errorf("fails to create user: %w", err)
	}

	return user, nil
}

/**
 * Helpers
 */

func createUser(user *User) error {
	if err := UserRepository.Save(user); err != nil {
		return fmt.Errorf("fails to create user: %w", err)
	}

	// create base todo list
	if _, err := createBasicTodoList(user); err != nil {
		return fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return nil
}

func createBasicTodoList(user *User) (TodoList, error) {
	basicTodoList := TodoList{
		// TODO[bug]: cycle dep
		// Name:    "Todos", // TODO i18n
		// IsBasic: true,
		// UserID:  user.ID,
	}
	if err := TodoListRepo.Save(&basicTodoList); err != nil {
		return TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	// TODO[bug]: cycle dep
	// user.BasicTodoListID = basicTodoList.ID
	if err := UserRepository.Save(user); err != nil {
		return TodoList{}, fmt.Errorf("fails to create user basic todo list: %w", err)
	}

	return basicTodoList, nil
}
