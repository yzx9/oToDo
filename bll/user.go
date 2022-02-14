package bll

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/otodo"
)

// User
type CreateUserPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func CreateUser(payload CreateUserPayload) (entity.User, error) {
	if len(payload.UserName) < 5 {
		return entity.User{}, fmt.Errorf("invalid user name: %v", payload.UserName)
	}

	if len(payload.Password) < 6 {
		return entity.User{}, fmt.Errorf("password too short")
	}

	basicTodoListID := uuid.NewString()
	user := entity.User{
		Entity: entity.Entity{
			ID: uuid.NewString(),
		},
		Name:            payload.UserName,
		Nickname:        payload.Nickname,
		Password:        GetCryptoPassword(payload.Password),
		BasicTodoListID: basicTodoListID,
	}
	err := dal.InsertUser(&user)

	// TODO handle error
	dal.InsertTodoList(&entity.TodoList{
		Entity: entity.Entity{
			ID: basicTodoListID,
		},
		Name:      "Todos", // TODO i18n
		Deletable: false,
		UserID:    user.ID,
	})

	return user, err
}

func GetUser(userID string) (entity.User, error) {
	return dal.SelectUser(userID)
}

// Invalid User Refresh Token

func CreateUserInvalidRefreshToken(userID, tokenID string) (entity.UserInvalidRefreshToken, error) {
	model := entity.UserInvalidRefreshToken{
		Entity: entity.Entity{
			ID:        uuid.NewString(),
			CreatedAt: time.Now(),
		},
		UserID:  userID,
		TokenID: tokenID,
	}
	if err := dal.InsertUserInvalidRefreshToken(model); err != nil {
		return entity.UserInvalidRefreshToken{}, fmt.Errorf("fails to make user refresh token invalid: %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID, tokenID string) (bool, error) {
	valid, err := dal.ExistUserInvalidRefreshToken(userID, tokenID)
	if err != nil {
		return false, fmt.Errorf("fails to get user refresh token: %w", err)
	}

	return valid, nil
}

// Password
func GetCryptoPassword(password string) [32]byte {
	pwd := sha256.Sum256(append([]byte(password), otodo.Conf.Secret.PasswordNonce...))
	return pwd
}
