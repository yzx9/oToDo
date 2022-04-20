package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yzx9/otodo/domain/identity"
)

const uriOAuthAuthorize = "https://github.com/login/oauth/authorize"
const uriOAuthAccessToken = "https://github.com/login/oauth/access_token"

var ServiceNotAvailable = fmt.Errorf("fails to fetch github access token")
var ServiceChanged = fmt.Errorf("github api seems to have changed")

func (g Adapter) CreateOAuthURI(state string) (string, error) {
	uri := uriOAuthAuthorize +
		"?client_id=" + g.config.ClientID +
		"&redirect_uri=" + g.config.OAuthRedirectURI +
		"&state=" + state
	return uri, nil
}

func (g Adapter) FetchOAuthToken(code string) (identity.ThirdPartyOAuthToken, error) {
	write := func(err error) (identity.ThirdPartyOAuthToken, error) {
		return identity.ThirdPartyOAuthToken{}, err
	}

	// Fetch access token
	vals := url.Values{}
	vals.Add("client_id", g.config.ClientID)
	vals.Add("client_secret", g.config.ClientSecret)
	vals.Add("code", code)
	vals.Add("redirect_uri", g.config.OAuthRedirectURI)

	req, err := http.NewRequest(http.MethodPost, uriOAuthAccessToken, strings.NewReader(vals.Encode()))
	if err != nil {
		return write(ServiceNotAvailable)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return write(ServiceNotAvailable)
	}

	if res.StatusCode != http.StatusOK {
		return write(ServiceNotAvailable)
	}

	// Parse access token
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return write(ServiceNotAvailable)
	}

	token := struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}{}
	if err := json.Unmarshal(body, &token); err != nil || token.TokenType != "bearer" {
		return write(ServiceChanged)
	}

	// TODO: auto map
	t := identity.NewThirdPartyOAuthToken(identity.ThirdPartyTokenTypeGithubAccessToken, identity.OAuthToken{
		AccessToken: token.AccessToken,
		Scope:       token.Scope,
		TokenType:   token.TokenType,
	})
	return t, nil
}
