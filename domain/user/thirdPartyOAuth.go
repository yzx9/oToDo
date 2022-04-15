package user

import (
	"fmt"
	"time"

	"github.com/yzx9/otodo/acl/github"
)

type ThirdPartyOAuthToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Active    bool
	Type      int8
	Token     string
	Scope     string

	UserID int64
}

func NewGithubOAuthToken(payload github.OAuthToken) ThirdPartyOAuthToken {
	token := ThirdPartyOAuthToken{
		Active: true,
		Type:   int8(ThirdPartyTokenTypeGithubAccessToken),
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

func (token *ThirdPartyOAuthToken) UpdateThirdPartyOAuthToken() error {
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
