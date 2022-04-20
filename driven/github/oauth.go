package github

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yzx9/otodo/domain/identity"
	"github.com/yzx9/otodo/util"
)

const uriOAuthAuthorize = "https://github.com/login/oauth/authorize"
const uriOAuthAccessToken = "https://github.com/login/oauth/access_token"

func (g Adapter) CreateOAuthURI(state string) (string, error) {
	uri := uriOAuthAuthorize +
		"?client_id=" + g.config.ClientID +
		"&redirect_uri=" + g.config.OAuthRedirectURI +
		"&state=" + state
	return uri, nil
}

func (g Adapter) FetchOAuthToken(code string) (identity.OAuthToken, error) {
	write := func(err error) (identity.OAuthToken, error) {
		return identity.OAuthToken{}, err
	}

	// Fetch access token
	vals := url.Values{}
	vals.Add("client_id", g.config.ClientID)
	vals.Add("client_secret", g.config.ClientSecret)
	vals.Add("code", code)
	vals.Add("redirect_uri", g.config.OAuthRedirectURI)

	req, err := http.NewRequest(http.MethodPost, uriOAuthAccessToken, strings.NewReader(vals.Encode()))
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

	token := struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}{}
	if err := json.Unmarshal(body, &token); err != nil || token.TokenType != "bearer" {
		// TODO[feat]: this is a fatal error as it usually means GitHub API changes
		return write(util.NewErrorWithUnknown("fails to parse github access token"))
	}

	// TODO: auto map
	return identity.OAuthToken{
		AccessToken: token.AccessToken,
		Scope:       token.Scope,
		TokenType:   token.TokenType,
	}, nil
}
