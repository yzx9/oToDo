package identity

/**
 * Config
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
 * Event publisher
 */

var EventPublisher interface {
	Publish(event string, payload []byte)
	Subscribe(event string, cb func([]byte)) func()
}

/**
 * Repository
 */

var UserRepository interface {
	Save(entity *User) error
	Find(id int64) (User, error)
	FindByUserName(username string) (User, error)
	FindByGithubID(githubID int64) (User, error)
	FindByTodo(todoID int64) (User, error)
	ExistByUserName(username string) (bool, error)
	ExistByGithubID(githubID int64) (bool, error)
}

var ThirdPartyOAuthTokenRepository interface {
	Save(entity *ThirdPartyOAuthToken) error
	SaveByUserIDAndType(entity *ThirdPartyOAuthToken) error
	ExistActiveOne(userID int64, tokenType ThirdPartyTokenType) (bool, error)
}

var UserInvalidRefreshTokenRepository interface {
	Save(entity *UserInvalidRefreshToken) error
	Exist(userID int64, tokenID string) (bool, error)
}

/**
 * GitHub
 */

var GithubAdapter interface {
	CreateOAuthURI(state string) (string, error)
	FetchOAuthToken(code string) (GithubOAuthToken, error)
	FetchUserPublicProfile(token string) (GithubUserPublicProfile, error)
}

type GithubOAuthToken struct {
	AccessToken string
	Scope       string
	TokenType   string
}

type GithubUserPublicProfile struct {
	ID        int64
	Name      string
	AvatarURL string
	Email     string
}
