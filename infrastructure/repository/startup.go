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

var db *gorm.DB
var NewID func() int64

func startUpDatabase() error {
	var err error
	write := func(err error) error {
		return util.NewError(errors.ErrDatabaseConnectFailed, "fails to connect database: %w", err)
	}

	// See https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	c := config.Database
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

func startUpIDGenerator() error {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("fails to create id generator")
	}

	NewID = func() int64 {
		return node.Generate().Int64()
	}

	return nil
}

func StartUp() error {
	return startUpDatabase()
}
