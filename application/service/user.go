package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/identity"
)

func CreateUser(newUser dto.NewUser) (identity.User, error) {
	return identity.RegisterUser(newUser.UserName, newUser.Nickname, newUser.Password)
}

func GetUser(userID int64) (identity.User, error) {
	u, err := UserRepository.Find(userID)
	if err != nil {
		return identity.User{}, fmt.Errorf("fails to get user: %w", err)
	}

	return u, nil
}
