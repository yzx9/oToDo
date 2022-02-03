package dal

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/yzx9/otodo/entity"
)

var users = make(map[uuid.UUID]entity.User)

func GetUser(id uuid.UUID) (entity.User, error) {
	user, ok := users[id]
	if !ok {
		return entity.User{}, fmt.Errorf("user not found: %v", id)
	}

	return user, nil
}

func GetUserByUserName(username string) (entity.User, error) {
	for _, user := range users {
		if user.Name == username {
			return user, nil
		}
	}

	return entity.User{}, fmt.Errorf("user not found: %v", username)
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

	return entity.User{}, errors.New("user not found")
}
