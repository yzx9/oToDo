package user

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/acl/github"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

const OAuthStateLen = 10

// TODO[perf]: redis
var oauthStates = make(map[string]time.Time)

func UpdateThirdPartyOAuthToken(token *repository.ThirdPartyOAuthToken) error {
	// TODO[bug]: handle error
	exist, err := repository.ThirdPartyOAuthTokenRepo.ExistActiveThirdPartyOAuthToken(token.UserID, repository.ThirdPartyTokenType(token.Type))
	if err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	handle := repository.ThirdPartyOAuthTokenRepo.UpdateThirdPartyOAuthToken
	if !exist {
		handle = repository.ThirdPartyOAuthTokenRepo.InsertThirdPartyOAuthToken
	}

	if err := handle(token); err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	return nil
}

func UpdateThirdPartyOAuthTokenAsync(token *repository.ThirdPartyOAuthToken) {
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

func FetchGithubOAuthToken(code, state string) (repository.ThirdPartyOAuthToken, error) {
	c := config.GitHub
	write := func(err error) (repository.ThirdPartyOAuthToken, error) {
		return repository.ThirdPartyOAuthToken{}, err
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

	return repository.ThirdPartyOAuthToken{
		Active: true,
		Type:   int8(repository.ThirdPartyTokenTypeGithubAccessToken),
		Token:  payload.AccessToken,
		Scope:  payload.Scope,
	}, nil
}
