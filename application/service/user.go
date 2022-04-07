package service

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
)

func GetUser(userID int64) (repository.User, error) {
	user, err := repository.UserRepo.Find(userID)
	if err != nil {
		return repository.User{}, fmt.Errorf("fails to get user: %w", err)
	}

	return user, nil
}
