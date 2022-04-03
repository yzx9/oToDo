package bll

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/infrastructure/util"
)

const tokenType = `bearer`
const authorizationRegexString = `^[Bb]earer (?P<token>[\w-]+.[\w-]+.[\w-]+)$`

var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func Login(payload dto.LoginDTO) (dto.SessionToken, error) {
	write := func() (dto.SessionToken, error) {
		return dto.SessionToken{}, util.NewErrorWithBadRequest("invalid credential")
	}

	user, err := repository.SelectUserByUserName(payload.UserName)
	if err != nil || user.Password == nil {
		return write()
	}

	if cryptoPwd := GetCryptoPassword(payload.Password); !bytes.Equal(user.Password, cryptoPwd) {
		return write()
	}

	if payload.RefreshTokenExpiresIn <= 0 {
		payload.RefreshTokenExpiresIn = config.Session.RefreshTokenExpiresInDefault
	} else if payload.RefreshTokenExpiresIn > config.Session.RefreshTokenExpiresInMax {
		payload.RefreshTokenExpiresIn = config.Session.RefreshTokenExpiresInMax
	}

	return newSessionToken(user, payload.RefreshTokenExpiresIn), nil
}

func LoginByGithubOAuth(code, state string) (dto.SessionToken, error) {
	token, err := FetchGithubOAuthToken(code, state)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to login: %w", err)
	}

	profile, err := FetchGithubUserPublicProfile(token.Token)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to fetch github user: %w", err)
	}

	user, err := getOrRegisterUserByGithub(profile)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to get user: %w", err)
	}

	go UpdateThirdPartyOAuthTokenAsync(&token)

	exp := config.Session.RefreshTokenExpiresInOAuth
	return newSessionToken(user, exp), nil
}

func Logout(userID int64, refreshTokenID string) error {
	_, err := CreateUserInvalidRefreshToken(userID, refreshTokenID)
	return err
}

func NewAccessToken(userID int64, refreshTokenID string) (dto.SessionToken, error) {
	user, err := repository.SelectUser(userID)
	if err != nil {
		return dto.SessionToken{}, fmt.Errorf("fails to get user, %w", err)
	}

	return newAccessToken(user, refreshTokenID), nil
}

func ParseSessionToken(token string) (*jwt.Token, error) {
	return ParseToken(token, &dto.SessionTokenClaims{})
}

func ParseAccessToken(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	token, err := ParseToken(matches[1], &dto.SessionTokenClaims{})
	if err != nil {
		return nil, fmt.Errorf("fails to parse access token: %w", err)
	}

	return token, nil
}

func ShouldRefreshAccessToken(oldAccessToken *jwt.Token) bool {
	if !oldAccessToken.Valid {
		return false
	}

	claims, ok := oldAccessToken.Claims.(*dto.SessionTokenClaims)
	if !ok || claims.ExpiresAt == 0 {
		return false
	}

	thd := config.Session.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

// generate access token only
func newAccessToken(user repository.User, refreshTokenID string) dto.SessionToken {
	exp := config.Session.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	claims := dto.SessionTokenClaims{
		TokenClaims:    NewClaims(user.ID, dur),
		RefreshTokenID: refreshTokenID,
	}
	token := NewToken(claims)

	return dto.SessionToken{
		AccessToken: token,
		TokenType:   tokenType,
		ExpiresIn:   int64(exp),
	}
}

// generate new access token and refresh token
func newSessionToken(user repository.User, refreshTokenExp int) dto.SessionToken {
	// refresh token
	dur := time.Duration(refreshTokenExp * int(time.Second))

	claims := dto.SessionTokenClaims{TokenClaims: NewClaims(user.ID, dur)}
	claims.Id = uuid.NewString()
	refreshToken := NewToken(claims)
	refreshTokenID := claims.Id

	// access token
	re := newAccessToken(user, refreshTokenID)

	re.RefreshToken = refreshToken
	return re
}
