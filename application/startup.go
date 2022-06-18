package application

import (
	"github.com/yzx9/otodo/adapter/driven/repository"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/file"
	"gorm.io/gorm"
)

func StatrUp(db *gorm.DB) error {
	service.UserRepository = repository.NewUserRepository(db)
	service.TodoRepository = repository.NewTodoRepository(db)
	service.TodoListRepository = repository.NewTodoListRepository(db)
	service.TodoListFolderRepository = repository.NewTodoListFolderRepository(db)
	service.SharingRepository = repository.NewSharingRepository(db)

	file.PermissionCheckerFactory.Register(file.FileTypeTodo, service.CanAccessTodoFile)

	return nil
}
