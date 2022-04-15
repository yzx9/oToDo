package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/user"
)

func CreateUser(newUser dto.NewUser) (user.User, error) {
	return user.CreateUser(newUser.ToEntity())
}

func GetUser(userID int64) (user.User, error) {
	u, err := UserRepository.Find(userID)
	if err != nil {
		return user.User{}, fmt.Errorf("fails to get user: %w", err)
	}

	return u, nil
}
