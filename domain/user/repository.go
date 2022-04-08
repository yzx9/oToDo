package user

import "github.com/yzx9/otodo/infrastructure/repository"

var UserRepo UserRepository
var UserInvalidRefreshTokenRepo UserInvalidRefreshTokenRepository
var ThirdPartyOAuthTokenRepo ThirdPartyOAuthTokenRepository
var TodoListRepo todoListRepository

type UserRepository interface {
	Save(user *repository.User) error

	Find(id int64) (repository.User, error)

	FindByUserName(username string) (repository.User, error)

	FindByGithubID(githubID int64) (repository.User, error)

	FindByTodo(todoID int64) (repository.User, error)

	ExistByUserName(username string) (bool, error)

	ExistByGithubID(githubID int64) (bool, error)
}

type UserInvalidRefreshTokenRepository interface {
	Save(entity *repository.UserInvalidRefreshToken) error

	Exist(userID int64, tokenID string) (bool, error)
}

type ThirdPartyOAuthTokenRepository interface {
	Save(entity *repository.ThirdPartyOAuthToken) error

	SaveByUserIDAndType(entity *repository.ThirdPartyOAuthToken) error

	ExistActiveOne(userID int64, tokenType repository.ThirdPartyTokenType) (bool, error)
}

type todoListRepository interface {
	Save(todoList *repository.TodoList) error
}
