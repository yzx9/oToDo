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

	err := dal.InitDatabase()
	if err != nil {
		return fmt.Errorf("fails to init database: %w", err)
	}

	return nil
}
