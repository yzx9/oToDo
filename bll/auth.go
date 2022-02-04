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
)

// Config
// TODO configurable
var passwordNonce = []byte("test_nonce")
var accessTokenExpiresIn = 15 * time.Minute
var refreshTokenExpiresIn = 15 * 24 * time.Hour
var accessTokenRefreshThreshold = 5 * time.Minute

// Constans
var accessTokenExpiresInSeconds = int64(accessTokenExpiresIn.Seconds())
var tokenType = "Bearer"
var authorizationRegexString = "^Bearer (?P<token>[\\w-]+.[\\w-]+.[\\w-]+)$"
var authorizationRegex = regexp.MustCompile(authorizationRegexString)

type AuthTokenResult struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
}

type AuthTokenClaims struct {
	jwt.StandardClaims

	UserID       string `json:"user_id"`
	UserNickname string `json:"user_nickname,omitempty"`
}

func Login(userName string, password string) (AuthTokenResult, error) {
	user, err := dal.GetUserByUserName(userName)
	if err != nil {
		return AuthTokenResult{}, fmt.Errorf("user not found: %v", userName)
	}

	pwd := sha256.Sum256(append([]byte(password), passwordNonce...))
	if !bytes.Equal(user.Password, pwd[:]) {
		return AuthTokenResult{}, fmt.Errorf("invalid credential")
	}

	return AuthTokenResult{
		AccessToken:  newAccessToken(user, accessTokenExpiresIn),
		TokenType:    tokenType,
		ExpiresIn:    accessTokenExpiresInSeconds,
		RefreshToken: newRefreshToken(user, refreshTokenExpiresIn),
	}, nil
}

func Logout(refreshToken *jwt.Token) {
	// TODO use jti for logout
}

func NewAccessToken(userID string) (AuthTokenResult, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return AuthTokenResult{}, fmt.Errorf("invalid id, %v", userID)
	}

	user, err := dal.GetUser(id)
	if err != nil {
		return AuthTokenResult{}, fmt.Errorf("fails to get user, %w", err)
	}

	token := newAccessToken(user, accessTokenExpiresIn)
	return AuthTokenResult{
		AccessToken: token,
		TokenType:   tokenType,
		ExpiresIn:   accessTokenExpiresInSeconds,
	}, nil
}

func ParseAuthToken(token string) (*jwt.Token, error) {
	return ParseToken(token, &AuthTokenClaims{})
}

func ParseAccessToken(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	return ParseToken(matches[1], &AuthTokenClaims{})
}

func ShouldRefreshAccessToken(oldAccessToken *jwt.Token) bool {
	if !oldAccessToken.Valid {
		return false
	}

	claims, ok := oldAccessToken.Claims.(*AuthTokenClaims)
	if !ok || claims.ExpiresAt == 0 {
		return false
	}

	return time.Now().Add(accessTokenRefreshThreshold).Unix() > claims.ExpiresAt
}

func newTokenClaims(userID uuid.UUID, exp time.Duration) AuthTokenClaims {
	return AuthTokenClaims{
		StandardClaims: NewClaims(exp),
		UserID:         userID.String(),
	}
}

func newAccessToken(user entity.User, exp time.Duration) string {
	claims := newTokenClaims(user.ID, exp)
	claims.UserNickname = user.Nickname
	return NewJwt(claims)
}

func newRefreshToken(user entity.User, exp time.Duration) string {
	claims := newTokenClaims(user.ID, exp)
	claims.Id = uuid.NewString()
	return NewJwt(claims)
}
