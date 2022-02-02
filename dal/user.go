package dal

import (
	"errors"

	"github.com/google/uuid"

	"github.com/yzx9/otodo/entity"
)

var users = make(map[uuid.UUID]entity.User)

func GetUser(id uuid.UUID) (entity.User, error) {
	user, ok := users[id]
	if !ok {
		return entity.User{}, errors.New("user not found")
	}

	return user, nil
}
