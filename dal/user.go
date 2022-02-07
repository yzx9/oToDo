package dal

import (
	"github.com/google/uuid"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var users = make(map[uuid.UUID]entity.User)
var invalidRefreshTokens = make(map[uuid.UUID]entity.UserRefreshToken)

func InsertUser(user entity.User) (entity.User, error) {
	users[user.ID] = user
	return user, nil
}

func GetUser(id uuid.UUID) (entity.User, error) {
	user, ok := users[id]
	if !ok {
		return entity.User{}, utils.NewErrorWithNotFound("user not found: %v", id)
	}

	return user, nil
}

func GetUserByUserName(username string) (entity.User, error) {
	for _, user := range users {
		if user.Name == username {
			return user, nil
		}
	}

	return entity.User{}, utils.NewErrorWithNotFound("user not found, username: %v", username)
}

func GetUserByTodo(todoID uuid.UUID) (entity.User, error) {
	todo, err := GetTodo(todoID)
	if err != nil {
		return entity.User{}, nil
	}

	for _, user := range users {
		if user.ID == todo.UserID {
			return user, nil
		}
	}

	return entity.User{}, utils.NewErrorWithNotFound("user not found, todo id: %v", todoID)
}

func InsertInvalidUserRefreshToken(entity entity.UserRefreshToken) (entity.UserRefreshToken, error) {
	invalidRefreshTokens[entity.ID] = entity
	return entity, nil
}

func ExistInvalidUserRefreshToken(userID uuid.UUID, tokenID uuid.UUID) bool {
	for _, token := range invalidRefreshTokens {
		if token.UserID == userID && token.TokenID == tokenID {
			return true
		}
	}

	return false
}
