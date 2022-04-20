package identity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/util"
)

const OAuthStateLen = 10

// TODO[perf]: redis
var oauthEntries = make(map[string]OAuth)

// entity
type OAuth struct {
	sessionID string
	state     string
	valid     bool
	createdAt time.Time
	expiresAt time.Time
}

func NewOAuthEntry() (OAuth, error) {
	now := time.Now()

	entry := OAuth{
		sessionID: uuid.NewString(),
		state:     util.RandomString(OAuthStateLen),
		valid:     true,
		createdAt: now,
		expiresAt: now.Add(time.Duration(Conf.OAuthStateExpiresIn * int(time.Second))),
	}

	oauthEntries[entry.state] = entry
	return entry, nil
}

func GetOAuthEntryByState(state string) (OAuth, error) {
	// Check state
	entry, ok := oauthEntries[state]
	if !ok || entry.expiresAt.Before(time.Now()) {
		// TODO: log
		return OAuth{}, InvalidCredential
	}

	delete(oauthEntries, entry.state)
	return entry, nil
}

func (a OAuth) GetGithubOAuthURI() (string, error) {
	uri, err := GithubAdapter.CreateOAuthURI(a.state)
	if err != nil {
		return "", newErr(fmt.Errorf("fails to create github oauth uri: %w", err))
	}

	return uri, nil
}

func (a OAuth) GetUserByGithub(code string) (User, error) {
	token, err := GithubAdapter.FetchOAuthToken(code)
	if err != nil {
		return User{}, newErr(fmt.Errorf("fails to fetch github oauth token: %w", err))
	}

	profile, err := GithubAdapter.FetchUserPublicProfile(token.Token)
	if err != nil {
		return User{}, newErr(fmt.Errorf("fails to fetch github user: %w", err))
	}

	u, err := GetOrRegisterUserByGithub(profile)
	if err != nil {
		return User{}, newErr(fmt.Errorf("fails to get user: %w", err))
	}

	return u, nil
}
