package identity

import (
	"fmt"
	"time"
)

// aggregate
type ThirdPartyOAuthToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	Type      ThirdPartyTokenType
	Token     string
	Scope     string

	UserID int64
}

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

type OAuthToken struct {
	AccessToken string
	Scope       string
	TokenType   string
}

func NewThirdPartyOAuthToken(tokenType ThirdPartyTokenType, payload OAuthToken) ThirdPartyOAuthToken {
	token := ThirdPartyOAuthToken{
		Active: true,
		Type:   ThirdPartyTokenTypeGithubAccessToken,
		Token:  payload.AccessToken,
		Scope:  payload.Scope,
	}

	go func() {
		if err := token.UpdateThirdPartyOAuthToken(); err != nil {
			// TODO[bug]: handle error
			fmt.Println(err)
		}
	}()

	return token
}

func (token *ThirdPartyOAuthToken) UpdateThirdPartyOAuthToken() (err error) {
	defer func() {
		if err != nil {
			newErr(fmt.Errorf("fails to update third party oauth token: %w", err))
		}
	}()

	// TODO[bug]: handle error
	exist, err := ThirdPartyOAuthTokenRepository.ExistActiveOne(token.UserID, ThirdPartyTokenType(token.Type))
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
