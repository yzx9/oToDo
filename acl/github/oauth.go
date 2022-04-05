package github

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/util"
)

const UriOAuthAuthorize = "https://github.com/login/oauth/authorize"
const UriOAuthAccessToken = "https://github.com/login/oauth/access_token"

func CreateOAuthURI(state string) (string, error) {
	c := config.GitHub
	uri := UriOAuthAuthorize
	uri += "?client_id=" + c.ClientID
	uri += "&redirect_uri=" + c.OAuthRedirectURI
	uri += "&state=" + state
	return uri, nil
}

type OAuthToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func FetchOAuthToken(code string) (OAuthToken, error) {
	c := config.GitHub
	write := func(err error) (OAuthToken, error) {
		return OAuthToken{}, err
	}

	// Fetch access token
	vals := url.Values{}
	vals.Add("client_id", c.ClientID)
	vals.Add("client_secret", c.ClientSecret)
	vals.Add("code", code)
	vals.Add("redirect_uri", c.OAuthRedirectURI)

	req, err := http.NewRequest(http.MethodPost, UriOAuthAccessToken, strings.NewReader(vals.Encode()))
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

	payload := OAuthToken{}
	if err := json.Unmarshal(body, &payload); err != nil || payload.TokenType != "bearer" {
		// TODO[feat]: this is a fatal error as it usually means GitHub API changes
		return write(util.NewErrorWithUnknown("fails to parse github access token"))
	}

	return payload, nil
}
