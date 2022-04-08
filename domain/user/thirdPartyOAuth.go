package user

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/acl/github"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/util"
)

const OAuthStateLen = 10

// TODO[perf]: redis
var oauthStates = make(map[string]time.Time)

type ThirdPartyOAuthToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	Active bool
	Type   int8
	Token  string
	Scope  string

	UserID int64
}

func UpdateThirdPartyOAuthToken(token *ThirdPartyOAuthToken) error {
	// TODO[bug]: handle error
	exist, err := ThirdPartyOAuthTokenRepository.ExistActiveOne(token.UserID, ThirdPartyTokenType(token.Type))
	if err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	save := ThirdPartyOAuthTokenRepository.SaveByUserIDAndType
	if !exist {
		save = ThirdPartyOAuthTokenRepository.Save
	}

	if err := save(token); err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	return nil
}

func UpdateThirdPartyOAuthTokenAsync(token *ThirdPartyOAuthToken) {
	if err := UpdateThirdPartyOAuthToken(token); err != nil {
		// TODO[bug]: handle error
		fmt.Println(err)
	}
}

func CreateGithubOAuthURI() (string, error) {
	state := util.RandomString(OAuthStateLen)
	oauthStates[state] = time.Now()
	uri, err := github.CreateOAuthURI(state)
	if err != nil {
		delete(oauthStates, state)
		return "", util.NewErrorWithUnknown("fails to create github oauth uri: %w", err)
	}

	return uri, nil
}

func FetchGithubOAuthToken(code, state string) (ThirdPartyOAuthToken, error) {
	c := config.GitHub
	write := func(err error) (ThirdPartyOAuthToken, error) {
		return ThirdPartyOAuthToken{}, err
	}

	// Check state
	createAt, ok := oauthStates[state]
	if !ok || createAt.Add(time.Duration(c.OAuthStateExpiresIn*int(time.Second))).Before(time.Now()) {
		// TODO[feat]: log
		return write(util.NewErrorWithForbidden("invalid state"))
	}
	delete(oauthStates, state)

	payload, err := github.FetchOAuthToken(code)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github oauth token"))
	}

	return ThirdPartyOAuthToken{
		Active: true,
		Type:   int8(ThirdPartyTokenTypeGithubAccessToken),
		Token:  payload.AccessToken,
		Scope:  payload.Scope,
	}, nil
}
