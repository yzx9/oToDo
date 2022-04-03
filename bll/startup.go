package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
)

var hasInit = false

func StartUp() error {
	if hasInit {
		return nil
	}

	hasInit = true

	if err := repository.StartUp(); err != nil {
		return fmt.Errorf("fails to init database: %w", err)
	}

	return nil
}
