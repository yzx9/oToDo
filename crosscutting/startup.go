package crosscutting

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure"
)

func StartUp() error {
	if err := infrastructure.StartUp(); err != nil {
		return fmt.Errorf("fails to start-up infrastructure: %w", err)
	}

	return nil
}
