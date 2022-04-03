package bll

import (
	"fmt"

	"github.com/yzx9/otodo/model/entity"
)

var hasInit = false

func StartUp() error {
	if hasInit {
		return nil
	}

	hasInit = true

	if err := entity.StartUp(); err != nil {
		return fmt.Errorf("fails to init database: %w", err)
	}

	return nil
}
