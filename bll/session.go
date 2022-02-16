package bll

import (
	"bytes"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/otodo"
)

// Constans
var tokenType = "Bearer"
var authorizationRegexString = "^[Bb]earer (?P<token>[\\w-]+.[\\w-]+.[\\w-]+)$"
var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func Login(userName, password string) (dto.SessionDTO, error) {
	user, err := dal.SelectUserByUserName(userName)
	if err != nil {
		return dto.SessionDTO{}, fmt.Errorf("user not found: %v", userName)
	}

	if cryptoPwd := GetCryptoPassword(password); !bytes.Equal(user.Password, cryptoPwd) {
		return dto.SessionDTO{}, fmt.Errorf("invalid credential")
	}

	refreshToken, refreshTokenID := newRefreshToken(user)
	re := newAccessTokenWithResult(user, refreshTokenID)
	re.RefreshToken = refreshToken
	return re, nil
}

func Logout(userID int64, refreshTokenID string) error {
	_, err := CreateUserInvalidRefreshToken(userID, refreshTokenID)
	return err
}

func NewAccessToken(userID int64, refreshTokenID string) (dto.SessionDTO, error) {
	user, err := dal.SelectUser(userID)
	if err != nil {
		return dto.SessionDTO{}, fmt.Errorf("fails to get user, %w", err)
	}

	return newAccessTokenWithResult(user, refreshTokenID), nil
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

	thd := otodo.Conf.Session.AccessTokenRefreshThreshold
	dur := time.Duration(thd * int(time.Second))
	return time.Now().Add(dur).Unix() > claims.ExpiresAt
}

func newAccessTokenWithResult(user entity.User, refreshTokenID string) dto.SessionDTO {
	exp := otodo.Conf.Session.AccessTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))
	return dto.SessionDTO{
		AccessToken: newAccessToken(user, refreshTokenID, dur),
		TokenType:   tokenType,
		ExpiresIn:   int64(exp),
	}
}

func newAccessToken(user entity.User, refreshTokenID string, exp time.Duration) string {
	claims := dto.SessionTokenClaims{
		TokenClaims:    NewClaims(user.ID, exp),
		UserNickname:   user.Nickname,
		RefreshTokenID: refreshTokenID,
	}
	return NewToken(claims)
}

func newRefreshToken(user entity.User) (string, string) {
	exp := otodo.Conf.Session.RefreshTokenExpiresIn
	dur := time.Duration(exp * int(time.Second))

	claims := dto.SessionTokenClaims{TokenClaims: NewClaims(user.ID, dur)}
	claims.Id = uuid.NewString()
	return NewToken(claims), claims.Id
}
