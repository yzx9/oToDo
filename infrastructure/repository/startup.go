package repository

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func StartUp() (*gorm.DB, error) {
	if err := startUpIDGenerator(); err != nil {
		return nil, err
	}

	db, err := startUpDatabase()
	if err != nil {
		return nil, err
	}

	if err := autoMigrate(db); err != nil {
		return nil, util.NewError(errors.ErrDatabaseConnectFailed, "fails to migrate database: %w", err)
	}

	return db, nil
}

func startUpIDGenerator() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("fails to create id generator")
	}

	newID = func() int64 {
		return node.Generate().Int64()
	}

	return nil
}

func startUpDatabase() (*gorm.DB, error) {
	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	c := config.Database
	dsn := fmt.Sprintf("%v:%v@%v(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", c.UserName, c.Password, c.Protocol, c.Host, c.Port, c.DatabaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, util.NewError(errors.ErrDatabaseConnectFailed, "fails to connect database: %w", err)
	}

	return db, nil
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&File{},

		&User{},
		&UserInvalidRefreshToken{},

		&Todo{},
		&TodoStep{},
		&TodoRepeatPlan{},

		&TodoList{},
		&TodoListFolder{},

		&Tag{},

		&Sharing{},

		&ThirdPartyOAuthToken{},
	)
}
