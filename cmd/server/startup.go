package main

import (
	"fmt"

	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/driven/github"
	"github.com/yzx9/otodo/infrastructure/repository"
	"gorm.io/gorm"
)

func startUp() error {
	db, err := repository.StartUp()
	if err != nil {
		return fmt.Errorf("fails to start-up infrastructure: %w", err)
	}

	if err := startUpDomain(db); err != nil {
		return fmt.Errorf("fails to start-up domain: %w", err)
	}

	if err := startUpApplication(db); err != nil {
		return fmt.Errorf("fails to start-up application: %w", err)
	}

	return nil
}

func startUpDomain(db *gorm.DB) error {
	if err := startUpFileDomain(db); err != nil {
		return err
	}

	if err := startUpIdentityDomain(db); err != nil {
		return err
	}

	if err := startUpTodoDomain(db); err != nil {
		return err
	}

	if err := startUpTodoListDomain(db); err != nil {
		return err
	}

	return nil
}

func startUpFileDomain(db *gorm.DB) error {
	file.FileRepository = repository.NewFileRepository(db)

	file.PermissionCheckerFactory.Register(file.FileTypeTodo, service.CanAccessTodoFile)

	return nil
}

func startUpIdentityDomain(db *gorm.DB) error {
	// config
	identity.Conf = config.IdentityDomain

	// repository
	identity.UserRepository = repository.NewUserRepository(db)
	identity.ThirdPartyOAuthTokenRepository = repository.NewThirdPartyOAuthTokenRepository(db)
	identity.UserInvalidRefreshTokenRepository = repository.NewUserInvalidRefreshTokenRepository(db)
	identity.TodoListRepo = repository.NewTodoListRepository(db)

	// driven
	identity.GithubAdapter = github.New(config.GitHubAdapter)

	return nil
}

func startUpTodoDomain(db *gorm.DB) error {
	todo.TodoRepository = repository.NewTodoRepository(db)
	todo.TodoStepRepository = repository.NewTodoStepRepository(db)
	todo.TodoRepeatPlanRepository = repository.NewTodoRepeatPlanRepository(db)
	todo.TagRepository = repository.NewTagRepository(db)
	todo.TagTodoRepository = repository.NewTagTodoRepository(db)
	todo.TodoFileRepository = repository.NewTodoFileRepository(db)
	return nil
}

func startUpTodoListDomain(db *gorm.DB) error {
	todolist.TodoRepository = repository.NewTodoRepository(db)
	todolist.TodoListRepository = repository.NewTodoListRepository(db)
	todolist.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	todolist.SharingRepository = repository.NewSharingRepository(db)
	todolist.TodoListSharingRepository = repository.NewTodoListSharingRepository(db)
	return nil
}

func startUpApplication(db *gorm.DB) error {
	service.UserRepository = repository.NewUserRepository(db)
	service.TodoRepository = repository.NewTodoRepository(db)
	service.TodoListRepository = repository.NewTodoListRepository(db)
	service.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	service.SharingRepository = repository.NewSharingRepository(db)

	return nil
}
