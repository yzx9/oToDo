package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
)

var hasInit = false

func Init() error {
	if hasInit {
		return nil
	}

	hasInit = true

	if err := dal.InitDatabase(); err != nil {
		return fmt.Errorf("fails to init database: %w", err)
	}

	return nil
}
