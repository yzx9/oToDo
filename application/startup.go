package application

import (
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/infrastructure"
)

func StartUp() error {
	if err := infrastructure.StartUp(); err != nil {
		return err
	}

	if err := bll.StartUp(); err != nil {
		return err
	}

	return nil
}
