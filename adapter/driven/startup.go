package driven

import (
	"fmt"

	"github.com/yzx9/otodo/adapter/driven/repository"
	"github.com/yzx9/otodo/infrastructure/event_publisher"
	infra_repo "github.com/yzx9/otodo/infrastructure/repository"
	"gorm.io/gorm"
)

func StartUp() (*gorm.DB, *event_publisher.EventPublisher, error) {
	db, err := infra_repo.StartUp()
	if err != nil {
		return nil, nil, fmt.Errorf("fails to start-up infrastructure repository: %w", err)
	}

	if err := autoMigrate(db); err != nil {
		return nil, nil, fmt.Errorf("fails to migrate database: %w", err)
	}

	ep, err := event_publisher.StartUp()
	if err != nil {
		return nil, nil, fmt.Errorf("fails to start-up infrastructure event publisher: %w", err)
	}

	return db, ep, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&repository.File{},

		&repository.User{},
		&repository.UserInvalidRefreshToken{},

		&repository.Todo{},
		&repository.TodoStep{},
		&repository.TodoRepeatPlan{},

		&repository.TodoList{},
		&repository.TodoListFolder{},

		&repository.Tag{},

		&repository.Sharing{},

		&repository.ThirdPartyOAuthToken{},
	)
}
