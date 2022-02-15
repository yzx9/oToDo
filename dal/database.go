package dal

import (
	"fmt"

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
	var err error

	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	c := otodo.Conf.Database
	dsn := fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.Protocol, c.Host, c.Port, c.DatabaseName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return write(err)
	}

	if err = autoMigrate(); err != nil {
		return write(err)
	}

	return nil
}

func autoMigrate() error {
	return db.AutoMigrate(
		&entity.File{},

		&entity.User{},
		&entity.UserInvalidRefreshToken{},

		&entity.Todo{},
		&entity.TodoFile{},
		&entity.TodoStep{},
		&entity.TodoRepeatPlan{},

		&entity.TodoList{},
		&entity.TodoListFolder{},

		&entity.Tag{},
	)
}
