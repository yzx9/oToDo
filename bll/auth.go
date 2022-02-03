package bll

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/utils"
)

// TODO configurable
var passwordNonce = []byte("test_nonce")
var accessTokenExpiresIn = 5 * time.Minute
var refreshTokenExpiresIn = 15 * 24 * time.Hour

type TokenResult struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
}

func Login(userName string, password string) (TokenResult, error) {
	user, err := dal.GetUserByUserName(userName)
	if err != nil {
		return TokenResult{}, fmt.Errorf("user not found: %v", userName)
	}

	pwd := sha256.Sum256(append([]byte(password), passwordNonce...))
	if !bytes.Equal(user.Password, pwd[:]) {
		return TokenResult{}, fmt.Errorf("invalid credential")
	}

	accessToken, exp := newAccessToken(user)
	return TokenResult{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		ExpiresIn:    exp,
		RefreshToken: newRefreshToken(user),
	}, nil
}

func Logout(refreshToken *jwt.Token) {
	// TODO use jti for logout
}

var authorizationRegexString = "^Bearer (?P<token>[\\w-]+.[\\w-]+.[\\w-]+)$"
var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func NewAccessToken(userID string) (TokenResult, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return TokenResult{}, fmt.Errorf("invalid id, %v", userID)
	}

	user, err := dal.GetUser(id)
	if err != nil {
		return TokenResult{}, fmt.Errorf("fails to get user, %w", err)
	}

	token, exp := newAccessToken(user)
	return TokenResult{
		AccessToken: token,
		TokenType:   "bearer",
		ExpiresIn:   exp,
	}, nil
}

func ParseAccessToken(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	return utils.ParseJWT(matches[1])
}

func newAccessToken(user entity.User) (string, int64) {
	claims := utils.NewTokenClaims(user.ID, accessTokenExpiresIn)
	claims.UserName = user.Name
	return utils.NewJwt(claims), claims.ExpiresAt
}

func newRefreshToken(user entity.User) string {
	claims := utils.NewTokenClaims(user.ID, refreshTokenExpiresIn)
	claims.Id = uuid.NewString()
	return utils.NewJwt(claims)
}
