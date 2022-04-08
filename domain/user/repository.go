package user

var UserRepository userRepository
var UserInvalidRefreshTokenRepository userInvalidRefreshTokenRepository
var ThirdPartyOAuthTokenRepository thirdPartyOAuthTokenRepository
var TodoListRepo todoListRepository

type userRepository interface {
	Save(user *User) error

	Find(id int64) (User, error)

	FindByUserName(username string) (User, error)

	FindByGithubID(githubID int64) (User, error)

	FindByTodo(todoID int64) (User, error)

	ExistByUserName(username string) (bool, error)

	ExistByGithubID(githubID int64) (bool, error)
}

type userInvalidRefreshTokenRepository interface {
	Save(entity *UserInvalidRefreshToken) error

	Exist(userID int64, tokenID string) (bool, error)
}

type thirdPartyOAuthTokenRepository interface {
	Save(entity *ThirdPartyOAuthToken) error

	SaveByUserIDAndType(entity *ThirdPartyOAuthToken) error

	ExistActiveOne(userID int64, tokenType ThirdPartyTokenType) (bool, error)
}

type todoListRepository interface {
	Save(todoList *TodoList) error
}

type TodoList struct{}
