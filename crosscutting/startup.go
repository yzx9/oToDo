package crosscutting

import (
	"fmt"

	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/domain/user"
	"github.com/yzx9/otodo/infrastructure/repository"
	"gorm.io/gorm"
)

func StartUp() error {
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
	file.FileRepository = repository.NewFileRepository(db)
	file.TodoFileRepository = repository.NewTodoFileRepository(db)

	todo.TodoRepository = repository.NewTodoRepository(db)
	todo.TodoStepRepository = repository.NewTodoStepRepository(db)
	todo.TodoRepeatPlanRepository = repository.NewTodoRepeatPlanRepository(db)
	todo.TagRepository = repository.NewTagRepository(db)
	todo.TagTodoRepository = repository.NewTagTodoRepository(db)

	todolist.TodoRepository = repository.NewTodoRepository(db)
	todolist.TodoListRepository = repository.NewTodoListRepository(db)
	todolist.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	todolist.SharingRepository = repository.NewSharingRepository(db)
	todolist.TodoListSharingRepository = repository.NewTodoListSharingRepository(db)

	user.UserRepository = repository.NewUserRepository(db)
	user.UserInvalidRefreshTokenRepository = repository.NewUserInvalidRefreshTokenRepository(db)
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
