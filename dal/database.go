package dal

import (
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() error {
	write := func(err error) error {
		return utils.NewError(otodo.ErrDatabaseConnectFailed, "fails to connect database: %w", err)
	}

	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := ""
	_db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return write(err)
	}
	db = _db

	if err = autoMigrate(); err != nil {
		return write(err)
	}

	return nil
}

func autoMigrate() error {
	return db.AutoMigrate(
		&entity.File{},

		&entity.User{},
		&entity.UserRefreshToken{},

		&entity.Todo{},
		&entity.TodoFile{},
		&entity.TodoStep{},
		&entity.TodoRepeatPlan{},

		&entity.TodoList{},
		&entity.TodoListFolder{},

		&entity.Tag{},
	)
}
