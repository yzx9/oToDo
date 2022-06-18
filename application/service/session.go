package service

import (
	"fmt"

	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/identity"
)

const TokenType = `bearer`

func Login(credential dto.UserCredential) (dto.SessionTokens, error) {
	user, err := identity.GetUserByUserName(credential.UserName)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	if !user.Password().Equals(credential.Password) {
		return dto.SessionTokens{}, identity.InvalidCredential
	}

	session := user.Session()
	accessToken, err := session.NewAccessToken()
	if err != nil {
		return dto.SessionTokens{}, err
	}

	refreshToken, err := session.NewRefreshToken(credential.RefreshTokenExpiresIn)
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
	user, err := identity.GetUserByGithubOAuth(code, state)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	session := user.Session()
	accessToken, err := session.NewAccessToken()
	if err != nil {
		return dto.SessionTokens{}, err
	}

	refreshToken, err := session.NewRefreshToken(identity.Conf.RefreshTokenExpiresInOAuth)
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
	user, err := identity.GetUserBySessionToken(identity.RefreshToken, token)
	if err != nil {
		return dto.SessionTokens{}, err
	}

	accessToken, err := user.Session().NewAccessToken()
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
	user, err := identity.GetUserBySessionToken(identity.AccessToken, token)
	if err != nil {
		return dto.SessionValidation{Valid: false}
	}

	session := user.Session()
	dto := dto.SessionValidation{
		Valid:          true,
		UserID:         session.UserID(),
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
	user, err := identity.GetUserBySessionToken(identity.AccessToken, accessToken)
	if err != nil {
		// TODO log
		fmt.Println(err.Error())
		return
	}

	if err := user.Session().Deactivate(); err != nil {
		// TODO log
		fmt.Println(err.Error())
		return
	}
}

func CreateGithubOAuthURI() (dto.OAuthRedirector, error) {
	uri, err := identity.NewGithubOAuthURI()
	if err != nil {
		return dto.OAuthRedirector{}, nil
	}

	return dto.OAuthRedirector{
		RedirectURI: uri,
	}, nil
}
