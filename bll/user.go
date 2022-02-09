package bll

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

// Config
// TODO configurable
var passwordNonce = []byte("test_nonce")

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

	basicTodoListID := uuid.New()
	user, err := dal.InsertUser(entity.User{
		ID:              uuid.New(),
		Name:            payload.UserName,
		Nickname:        payload.Nickname,
		Password:        GetCryptoPassword(payload.Password),
		BasicTodoListID: basicTodoListID,
	})

	dal.InsertTodoList(entity.TodoList{
		ID:        basicTodoListID,
		Name:      "Todos", // TODO i18n
		Deletable: false,
		UserID:    user.ID,
	})

	return user, err
}

func GetUser(userID uuid.UUID) (entity.User, error) {
	return dal.GetUser(userID)
}

// Invalid User Refresh Token

func CreateInvalidUserRefreshToken(userID, tokenID uuid.UUID) (entity.UserRefreshToken, error) {
	model, err := dal.InsertInvalidUserRefreshToken(entity.UserRefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenID:   tokenID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("fails to invalid user refresh token, %w", err)
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID, tokenID uuid.UUID) bool {
	return !dal.ExistInvalidUserRefreshToken(userID, tokenID)
}

// Password
func GetCryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), passwordNonce...))
	return pwd[:]
}
