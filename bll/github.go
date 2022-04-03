package bll

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const githubOAuthStateLen = 10
const githubOAuthAuthorizeURI = "https://github.com/login/oauth/authorize"
const githubOAuthAccessTokenURI = "https://github.com/login/oauth/access_token"
const githubUserURI = "https://api.github.com/user"

// TODO[perf]: redis
var githubOAuthStates = make(map[string]time.Time)

func CreateGithubOAuthURI() (string, error) {
	c := config.GitHub
	uri := githubOAuthAuthorizeURI

	uri += "?client_id=" + c.ClientID
	uri += "&redirect_uri=" + c.OAuthRedirectURI

	state := util.RandomString(githubOAuthStateLen)
	githubOAuthStates[state] = time.Now()
	uri += "&state=" + state

	return uri, nil
}

func FetchGithubOAuthToken(code, state string) (entity.ThirdPartyOAuthToken, error) {
	c := config.GitHub
	write := func(err error) (entity.ThirdPartyOAuthToken, error) {
		return entity.ThirdPartyOAuthToken{}, err
	}

	// Check state
	createAt, ok := githubOAuthStates[state]
	if !ok || createAt.Add(time.Duration(c.OAuthStateExpiresIn*int(time.Second))).Before(time.Now()) {
		// TODO[feat]: log
		return write(util.NewErrorWithForbidden("invalid state"))
	}
	delete(githubOAuthStates, state)

	// Fetch access token
	vals := url.Values{}
	vals.Add("client_id", c.ClientID)
	vals.Add("client_secret", c.ClientSecret)
	vals.Add("code", code)
	vals.Add("redirect_uri", c.OAuthRedirectURI)

	req, err := http.NewRequest(http.MethodPost, githubOAuthAccessTokenURI, strings.NewReader(vals.Encode()))
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to new request"))
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github access token"))
	}

	if res.StatusCode != http.StatusOK {
		return write(util.NewErrorWithForbidden("invalid code"))
	}

	// Parse access token
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github access token"))
	}

	payload := dto.GithubOAuthAccessToken{}
	if err := json.Unmarshal(body, &payload); err != nil || payload.TokenType != "bearer" {
		// TODO[feat]: this is a fatal error as it usually means GitHub API changes
		return write(util.NewErrorWithUnknown("fails to parse github access token"))
	}

	return entity.ThirdPartyOAuthToken{
		Active: true,
		Type:   int8(entity.ThirdPartyTokenTypeGithubAccessToken),
		Token:  payload.AccessToken,
		Scope:  payload.Scope,
	}, nil
}

func FetchGithubUserPublicProfile(token string) (dto.GithubUserPublicProfile, error) {
	write := func(err error) (dto.GithubUserPublicProfile, error) {
		return dto.GithubUserPublicProfile{}, err
	}

	req, err := http.NewRequest(http.MethodGet, githubUserURI, strings.NewReader(""))
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to new request"))
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github user public profile"))
	}

	if res.StatusCode == http.StatusForbidden {
		// TODO[feat]: should we inactive token?
		return write(util.NewError(otodo.ErrThirdPartyForbidden, "github access token has been invalid"))
	}

	if res.StatusCode != http.StatusOK {
		return write(util.NewError(otodo.ErrThirdPartyUnknown, "fails to fetch github user public profile"))
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return write(util.NewErrorWithUnknown("fails to fetch github user public profile"))
	}

	payload := dto.GithubUserPublicProfile{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return write(util.NewErrorWithUnknown("fails to parse github user public profile"))
	}

	return payload, nil
}
