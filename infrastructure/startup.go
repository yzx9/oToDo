package infrastructure

import "github.com/yzx9/otodo/infrastructure/repository"

func StartUp() error {
	return repository.StartUp()
}
