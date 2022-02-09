package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

var users = make(map[string]entity.User)
var invalidRefreshTokens = make(map[string]entity.UserRefreshToken)

func InsertUser(user entity.User) error {
	users[user.ID] = user
	return nil
}

func GetUser(id string) (entity.User, error) {
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

func GetUserByTodo(todoID string) (entity.User, error) {
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

func InsertInvalidUserRefreshToken(entity entity.UserRefreshToken) error {
	invalidRefreshTokens[entity.ID] = entity
	return nil
}

func ExistInvalidUserRefreshToken(userID, tokenID string) bool {
	for _, token := range invalidRefreshTokens {
		if token.UserID == userID && token.TokenID == tokenID {
			return true
		}
	}

	return false
}
