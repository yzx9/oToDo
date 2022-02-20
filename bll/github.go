package bll

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const githubOAuthStateLen = 10

// TODO[perf]: redis
var githubOAuthStates = make(map[string]time.Time)

func CreateGithubOAuth() (string, error) {
	c := otodo.Conf.Github
	uri := c.OAuthAuthorizeURI

	uri += "?client_id=" + c.ClientID
	uri += "&redirect_uri=" + c.OAuthRedirectURI

	state := util.RandomString(githubOAuthStateLen)
	githubOAuthStates[state] = time.Now()
	uri += "&state=" + state

	return uri, nil
}

func FetchGithubOAuthToken(code, state string) (entity.ThirdPartyToken, error) {
	c := otodo.Conf.Github
	write := func(err error) (entity.ThirdPartyToken, error) {
		return entity.ThirdPartyToken{}, err
	}

	// Check state
	createAt, ok := githubOAuthStates[state]
	if !ok || createAt.Add(time.Duration(c.OAuthStateExpiresIn*int(time.Second))).After(time.Now()) {
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

	req, err := http.NewRequest("POST", c.OAuthAccessTokenURI, strings.NewReader(vals.Encode()))
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

	type Payload struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}
	payload := Payload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return write(util.NewErrorWithUnknown("fails to parse github access token"))
	}

	// Save access token
	token := entity.ThirdPartyToken{
		Active:    true,
		Type:      int8(entity.ThirdPartyTokenTypeGithubAccessToken),
		Token:     payload.AccessToken,
		TokenType: payload.TokenType,
		Scope:     payload.Scope,
	}
	if err := dal.InsertThirdPartyToken(&token); err != nil {
		return write(fmt.Errorf("fails to save token: %w", err))
	}

	return token, nil
}
