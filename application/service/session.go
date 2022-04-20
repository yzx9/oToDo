package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/config"
	"github.com/yzx9/otodo/domain/session"
)

const TokenType = `bearer`

func Login(credential dto.UserCredential) (dto.SessionTokens, error) {
	tokens, err := session.LoginByCredential(credential.UserName, credential.Password)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	accessToken, err := tokens.NewAccessToken()
	if err != nil {
		return dto.SessionTokens{}, err
	}

	refreshToken, err := tokens.NewRefreshToken(credential.RefreshTokenExpiresIn)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	return dto.SessionTokens{
		AccessToken:  accessToken.Token,
		TokenType:    TokenType,
		ExpiresIn:    accessToken.ExpiresIn,
		RefreshToken: refreshToken.Token,
	}, nil
}

func LoginByGithubOAuth(code string, state string) (dto.SessionTokens, error) {
	session, err := session.LoginByGithubOAuth(code, state)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	accessToken, err := session.NewAccessToken()
	if err != nil {
		return dto.SessionTokens{}, err
	}

	refreshToken, err := session.NewRefreshToken(config.Session.RefreshTokenExpiresInOAuth)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	return dto.SessionTokens{
		AccessToken:  accessToken.Token,
		TokenType:    TokenType,
		ExpiresIn:    accessToken.ExpiresIn,
		RefreshToken: refreshToken.Token,
	}, nil
}

func LoginByRefreshToken(token string) (dto.SessionTokens, error) {
	session, err := session.LoginByRefreshToken(token)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	accessToken, err := session.NewAccessToken()
	if err != nil {
		return dto.SessionTokens{}, err
	}

	return dto.SessionTokens{
		AccessToken:  accessToken.Token,
		TokenType:    TokenType,
		ExpiresIn:    accessToken.ExpiresIn,
		RefreshToken: token,
	}, nil
}

func LoginByAccessToken(token string) dto.SessionValidation {
	session, err := session.LoginByAccessToken(token)
	if err != nil {
		return dto.SessionValidation{Valid: false}
	}

	dto := dto.SessionValidation{
		Valid:          true,
		UserID:         session.UserID,
		NewAccessToken: false,
	}

	if !session.ShouldRefreshAccessToken(token) {
		return dto
	}

	accessToken, err := session.NewAccessToken()
	if err != nil {
		return dto
	}

	dto.NewAccessToken = true
	dto.AccessToken = accessToken.Token
	return dto
}

func Logout(accessToken string) {
	session, err := session.LoginByAccessToken(accessToken)
	if err != nil {
		// TODO log
		fmt.Println(err.Error())
		return
	}

	if err := session.Logout(); err != nil {
		// TODO log
		fmt.Println(err.Error())
		return
	}
}

func CreateGithubOAuthURI() (dto.OAuthRedirector, error) {
	oauth, err := session.NewOAuthEntry()
	if err != nil {
		return dto.OAuthRedirector{}, nil
	}

	uri, err := oauth.GetGithubOAuthURI()
	if err != nil {
		return dto.OAuthRedirector{}, nil
	}

	return dto.OAuthRedirector{
		RedirectURI: uri,
	}, nil
}
