package identity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/util"
)

const OAuthStateLen = 10

// TODO[perf]: redis
var oauthEntries = make(map[string]oauth)

// entity
type oauth struct {
	sessionID string
	state     string
	valid     bool
	createdAt time.Time
	expiresAt time.Time
}

func newOAuthEntry() (oauth, error) {
	now := time.Now()
	entry := oauth{
		sessionID: uuid.NewString(),
		state:     util.RandomString(OAuthStateLen),
		valid:     true,
		createdAt: now,
		expiresAt: now.Add(time.Duration(Conf.OAuthStateExpiresIn * int(time.Second))),
	}

	oauthEntries[entry.state] = entry
	return entry, nil
}

func getOAuthEntryByState(state string) (oauth, error) {
	// Check state
	entry, ok := oauthEntries[state]
	if !ok || entry.expiresAt.Before(time.Now()) {
		// TODO: log
		return oauth{}, InvalidCredential
	}

	delete(oauthEntries, entry.state)
	return entry, nil
}

func (a oauth) GetGithubOAuthURI() (string, error) {
	uri, err := GithubAdapter.CreateOAuthURI(a.state)
	if err != nil {
		return "", Error{fmt.Errorf("fails to create github oauth uri: %w", err)}
	}

	return uri, nil
}

func (a oauth) GetOrRegisterUserByGithub(code string) (User, error) {
	githubToken, err := GithubAdapter.FetchOAuthToken(code)
	if err != nil {
		return User{}, Error{fmt.Errorf("fails to fetch github oauth token: %w", err)}
	}

	token := ThirdPartyOAuthToken{
		active: true,
		typee:  ThirdPartyTokenTypeGithubAccessToken,
		token:  githubToken.AccessToken,
		scope:  githubToken.Scope,
	}

	go func() {
		if err := token.Save(); err != nil {
			// TODO: handle error
			fmt.Println(err)
		}
	}()

	profile, err := GithubAdapter.FetchUserPublicProfile(token.token)
	if err != nil {
		return User{}, Error{fmt.Errorf("fails to fetch github user: %w", err)}
	}

	exist, err := UserRepository.ExistByGithubID(profile.ID)
	if err != nil {
		return User{}, Error{fmt.Errorf("fails to register user: %w", err)}
	}

	if exist {
		user, err := UserRepository.FindByGithubID(profile.ID)
		if err != nil {
			return User{}, Error{fmt.Errorf("fails to get user: %w", err)}
		}

		return user, nil
	}

	// Register new user
	// TODO[feat]: download user avatar
	user := User{
		name:     profile.Email,
		nickname: profile.Name,
		email:    profile.Email,
		githubId: profile.ID,
	}
	if err := user.new(); err != nil {
		return User{}, Error{fmt.Errorf("fails to get user: %w", err)}
	}

	return user, nil
}

// aggregate
type ThirdPartyOAuthToken struct {
	id        int64
	createdAt time.Time
	updatedAt time.Time
	userId    int64
	active    bool
	typee     ThirdPartyTokenType
	token     string
	scope     string
}

// value object
type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

func NewThirdPartyOAuthToken(
	id int64,
	createdAt time.Time,
	updatedAt time.Time,
	userId int64,
	active bool,
	typee ThirdPartyTokenType,
	token string,
	scope string,
) ThirdPartyOAuthToken {
	return ThirdPartyOAuthToken{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
		userId:    userId,
		active:    active,
		typee:     typee,
		token:     token,
		scope:     scope,
	}
}

func (token ThirdPartyOAuthToken) ID() int64                 { return token.id }
func (token ThirdPartyOAuthToken) CreatedAt() time.Time      { return token.createdAt }
func (token ThirdPartyOAuthToken) UpdatedAt() time.Time      { return token.updatedAt }
func (token ThirdPartyOAuthToken) UserID() int64             { return token.userId }
func (token ThirdPartyOAuthToken) Active() bool              { return token.active }
func (token ThirdPartyOAuthToken) Type() ThirdPartyTokenType { return token.typee }
func (token ThirdPartyOAuthToken) Token() string             { return token.token }
func (token ThirdPartyOAuthToken) Scope() string             { return token.scope }

func (token ThirdPartyOAuthToken) SetID(id int64) {
	if token.id == 0 {
		token.id = id
	}
}

func (token *ThirdPartyOAuthToken) Save() (err error) {
	defer func() {
		if err != nil {
			err = Error{fmt.Errorf("fails to update third party oauth token: %w", err)}
		}
	}()

	// TODO[bug]: handle error
	exist, err := ThirdPartyOAuthTokenRepository.ExistActiveOne(token.userId, ThirdPartyTokenType(token.typee))
	if err != nil {
		return
	}

	save := ThirdPartyOAuthTokenRepository.SaveByUserIDAndType
	if !exist {
		save = ThirdPartyOAuthTokenRepository.Save
	}

	if err = save(token); err != nil {
		return
	}

	return nil
}
