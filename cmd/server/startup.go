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
	// File
	file.FileRepository = repository.NewFileRepository(db)
	file.PermissionCheckerFactory.Register(file.FileTypeTodo, service.CanAccessTodoFile)

	if err := startUpIdentityDomain(db); err != nil {
		return nil
	}

	// Todo
	todo.TodoRepository = repository.NewTodoRepository(db)
	todo.TodoStepRepository = repository.NewTodoStepRepository(db)
	todo.TodoRepeatPlanRepository = repository.NewTodoRepeatPlanRepository(db)
	todo.TagRepository = repository.NewTagRepository(db)
	todo.TagTodoRepository = repository.NewTagTodoRepository(db)
	todo.TodoFileRepository = repository.NewTodoFileRepository(db)

	// Todo List
	todolist.TodoRepository = repository.NewTodoRepository(db)
	todolist.TodoListRepository = repository.NewTodoListRepository(db)
	todolist.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	todolist.SharingRepository = repository.NewSharingRepository(db)
	todolist.TodoListSharingRepository = repository.NewTodoListSharingRepository(db)

	return nil
}

func startUpIdentityDomain(db *gorm.DB) error {
	identity.UserRepository = repository.NewUserRepository(db)
	identity.ThirdPartyOAuthTokenRepository = repository.NewThirdPartyOAuthTokenRepository(db)
	identity.UserInvalidRefreshTokenRepository = repository.NewUserInvalidRefreshTokenRepository(db)
	identity.TodoListRepo = repository.NewTodoListRepository(db)

	// TODO: auto map
	identity.GithubAdapter = github.New(github.Config{
		ClientID:            config.GitHub.ClientID,
		ClientSecret:        config.GitHub.ClientSecret,
		OAuthRedirectURI:    config.GitHub.OAuthRedirectURI,
		OAuthStateExpiresIn: config.GitHub.OAuthStateExpiresIn,
	})

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
