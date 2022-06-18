package domain

import (
	"github.com/yzx9/otodo/adapter/driven/github"
	"github.com/yzx9/otodo/adapter/driven/repository"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/domain/sharing"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/infrastructure/event_publisher"
	"gorm.io/gorm"
)

func StartUp(db *gorm.DB, ep *event_publisher.EventPublisher) error {
	if err := startUpFileDomain(db); err != nil {
		return err
	}

	if err := startUpIdentityDomain(db, ep); err != nil {
		return err
	}

	if err := startUpTodoDomain(db, ep); err != nil {
		return err
	}

	if err := startUpSharingDomain(db); err != nil {
		return err
	}

	return nil
}

func startUpFileDomain(db *gorm.DB) error {
	file.FileRepository = repository.NewFileRepository(db)

	return nil
}

func startUpIdentityDomain(db *gorm.DB, ep *event_publisher.EventPublisher) error {
	// config
	identity.Conf = config.IdentityDomain

	// event
	identity.EventPublisher = ep

	// repository
	identity.UserRepository = repository.NewUserRepository(db)
	identity.ThirdPartyOAuthTokenRepository = repository.NewThirdPartyOAuthTokenRepository(db)
	identity.UserInvalidRefreshTokenRepository = repository.NewUserInvalidRefreshTokenRepository(db)

	// driven
	identity.GithubAdapter = github.New(config.GitHubAdapter)

	return nil
}

func startUpTodoDomain(db *gorm.DB, ep *event_publisher.EventPublisher) error {
	// event
	todo.EventPublisher = ep
	ep.Subscribe(identity.EventUserCreated, todo.HandleUserCreatedEvent)

	todo.TodoRepository = repository.NewTodoRepository(db)
	todo.TodoStepRepository = repository.NewTodoStepRepository(db)
	todo.TodoRepeatPlanRepository = repository.NewTodoRepeatPlanRepository(db)
	todo.TagRepository = repository.NewTagRepository(db)
	todo.TagTodoRepository = repository.NewTagTodoRepository(db)
	todo.TodoFileRepository = repository.NewTodoFileRepository(db)

	// repository
	todo.TodoRepository = repository.NewTodoRepository(db)
	todo.TodoListRepository = repository.NewTodoListRepository(db)
	todo.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	todo.TodoListSharingRepository = repository.NewTodoListSharingRepository(db)

	return nil
}

func startUpSharingDomain(db *gorm.DB) error {
	sharing.SharingRepository = repository.NewSharingRepository(db)
	return nil
}
