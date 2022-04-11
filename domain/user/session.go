package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/acl/github"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

type Session struct {
	SessionID string
	UserID    int64
}

type SessionTokenClaims struct {
	JWTClaims

	RefreshTokenID string `json:"rti,omitempty"`
}

type Token struct {
	Token     string
	Type      TokenType
	ExpiresIn int64
}

func LoginByCredential(UserName, password string) (Session, error) {
	user, err := UserRepository.FindByUserName(UserName)
	if err != nil ||
		user.Password == nil ||
		!user.IsSamePassword(password) {
		return Session{}, util.NewErrorWithBadRequest("invalid credential")
	}

	return Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}, nil
}

func LoginByGithubOAuth(code, state string) (Session, error) {
	token, err := FetchGithubOAuthToken(code, state)
	if err != nil {
		return Session{}, fmt.Errorf("fails to login: %w", err)
	}

	profile, err := github.FetchGithubUserPublicProfile(token.Token)
	if err != nil {
		return Session{}, fmt.Errorf("fails to fetch github user: %w", err)
	}

	user, err := getOrRegisterUserByGithub(profile)
	if err != nil {
		return Session{}, fmt.Errorf("fails to get user: %w", err)
	}

	go UpdateThirdPartyOAuthTokenAsync(&token)

	return Session{
		SessionID: uuid.NewString(),
		UserID:    user.ID,
	}, nil
}

func LoginByAccessToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, util.NewError(errors.ErrUnauthorized, "invalid token")
	}

	return Session{
		SessionID: claims.RefreshTokenID,
		UserID:    claims.UserID,
	}, nil
}

func LoginByRefreshToken(token string) (Session, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return Session{}, util.NewError(errors.ErrUnauthorized, "invalid token")
	}

	valid, err := UserInvalidRefreshTokenRepository.Exist(claims.UserID, claims.Id)
	if err != nil || !valid {
		return Session{}, fmt.Errorf("invalid refresh token: %w", err)
	}

	return Session{
		SessionID: claims.Id,
		UserID:    claims.UserID,
	}, nil
}

func (s Session) Logout() error {
	_, err := CreateUserInvalidRefreshToken(s.UserID, s.SessionID)
	return err
}

func (s Session) NewAccessToken() (Token, error) {
	user, err := UserRepository.Find(s.UserID)
	if err != nil {
		return Token{}, fmt.Errorf("fails to get user, %w", err)
	}

	exp := config.Session.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	claims := SessionTokenClaims{
		JWTClaims:      NewClaims(user.ID, dur),
		RefreshTokenID: s.SessionID,
	}
	token := NewToken(claims)

	return Token{
		Token:     token,
		Type:      AccessToken,
		ExpiresIn: int64(exp),
	}, nil
}

// generate refresh token
func (s Session) NewRefreshToken(exp int) Token {
	if exp <= 0 {
		exp = config.Session.RefreshTokenExpiresInDefault
	} else if exp > config.Session.RefreshTokenExpiresInMax {
		exp = config.Session.RefreshTokenExpiresInMax
	}
	dur := time.Duration(exp * int(time.Second))

	claims := SessionTokenClaims{JWTClaims: NewClaims(s.UserID, dur)}
	claims.Id = uuid.NewString()
	token := NewToken(claims)
	return Token{
		Token:     token,
		Type:      RefreshToken,
		ExpiresIn: int64(exp),
	}
}

func (s Session) ShouldRefreshAccessToken(accessToken string) bool {
	claims, err := parseSessionToken(accessToken)
	if err != nil {
		return false
	}

	thd := config.Session.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

func parseSessionToken(tokenString string) (*SessionTokenClaims, error) {
	token, err := ParseToken(tokenString, &SessionTokenClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SessionTokenClaims)
	if !ok {
		return nil, fmt.Errorf("unknown error")
	}

	return claims, nil
}
