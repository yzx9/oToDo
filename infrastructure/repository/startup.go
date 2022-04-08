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

func StartUpDatabase() (*gorm.DB, error) {
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

func startUpRepositories(db *gorm.DB) {
	UserRepo = UserRepository{db: db}
	UserInvalidRefreshTokenRepo = UserInvalidRefreshTokenRepository{db: db}

	TodoRepo = TodoRepository{db: db}
	TodoStepRepo = TodoStepRepository{db: db}
	TodoRepeatPlanRepo = TodoRepeatPlanRepository{db: db}

	TodoListRepo = TodoListRepository{db: db}
	TodoListSharingRepo = TodoListSharingRepository{db: db}
	TodoListFolderRepo = TodoListFolderRepository{db: db}

	TagRepo = TagRepository{db: db}
	TagTodoRepo = TagTodoRepository{db: db}

	SharingRepo = SharingRepository{db: db}

	ThirdPartyOAuthTokenRepo = ThirdPartyOAuthTokenRepository{db: db}
}

func StartUp() (*gorm.DB, error) {
	if err := startUpIDGenerator(); err != nil {
		return nil, err
	}

	db, err := StartUpDatabase()
	if err != nil {
		return nil, err
	}

	if err := autoMigrate(db); err != nil {
		return nil, util.NewError(errors.ErrDatabaseConnectFailed, "fails to migrate database: %w", err)
	}

	startUpRepositories(db)

	return db, nil
}
