package identity

import (
	"github.com/yzx9/otodo/domain/todolist"
)

/**
 * config
 */
var Conf Config

type Config struct {
	AccessTokenExpiresIn         int // second
	AccessTokenRefreshThreshold  int // second
	RefreshTokenExpiresInDefault int // second
	RefreshTokenExpiresInMax     int // second
	RefreshTokenExpiresInOAuth   int // second
	OAuthStateExpiresIn          int // second
	TokenIssuer                  string

	// secret
	PasswordNonce   []byte
	TokenHmacSecret []byte
}

/**
 * repository
 */

var UserRepository userRepository
var ThirdPartyOAuthTokenRepository thirdPartyOAuthTokenRepository
var UserInvalidRefreshTokenRepository userInvalidRefreshTokenRepository
var TodoListRepo todoListRepository

type userRepository interface {
	Save(entity *User) error
	Find(id int64) (User, error)
	FindByUserName(username string) (User, error)
	FindByGithubID(githubID int64) (User, error)
	FindByTodo(todoID int64) (User, error)
	ExistByUserName(username string) (bool, error)
	ExistByGithubID(githubID int64) (bool, error)
}

type thirdPartyOAuthTokenRepository interface {
	Save(entity *ThirdPartyOAuthToken) error
	SaveByUserIDAndType(entity *ThirdPartyOAuthToken) error
	ExistActiveOne(userID int64, tokenType ThirdPartyTokenType) (bool, error)
}

type userInvalidRefreshTokenRepository interface {
	Save(entity *UserInvalidRefreshToken) error
	Exist(userID int64, tokenID string) (bool, error)
}

type todoListRepository interface {
	Save(entity *todolist.TodoList) error
}

/**
 * GitHub
 */

var GithubAdapter githubAdapter

type githubAdapter interface {
	CreateOAuthURI(state string) (string, error)
	FetchOAuthToken(code string) (ThirdPartyOAuthToken, error)
	FetchUserPublicProfile(token string) (GithubUserPublicProfile, error)
}

type GithubUserPublicProfile struct {
	ID        int64
	Name      string
	AvatarURL string
	Email     string
}
