package main

import (
	"fmt"

	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/session"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/domain/user"
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

	// Session
	session.UserRepository = repository.NewUserRepository(db)
	session.UserInvalidRefreshTokenRepository = repository.NewUserInvalidRefreshTokenRepository(db)

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

	// User
	user.UserRepository = repository.NewUserRepository(db)
	user.ThirdPartyOAuthTokenRepository = repository.NewThirdPartyOAuthTokenRepository(db)
	user.TodoListRepo = repository.NewTodoListRepository(db)

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
