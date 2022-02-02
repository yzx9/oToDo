package bll

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

// TODO configurable
var passwordNonce = []byte("test_nonce")
var tokenIssuer = "oToDo"
var tokenHmacSecret = []byte("test_secret")
var accessTokenExpiresIn = 5 * time.Minute
var refreshTokenExpiresIn = 15 * 24 * time.Hour

type LoginResult struct {
	AccessToken  string
	TokenType    string
	ExpiresIn    int64
	RefreshToken string
}

func Login(userName string, password string) (LoginResult, error) {
	user, err := dal.GetUserByUserName(userName)
	if err != nil {
		return LoginResult{}, fmt.Errorf("user not found: %v", userName)
	}

	pwd := sha256.Sum256(append([]byte(password), passwordNonce...))
	fmt.Println(hex.EncodeToString(pwd[:]))
	if !bytes.Equal(user.Password, pwd[:]) {
		return LoginResult{}, fmt.Errorf("invalid credential")
	}

	accessToken, exp := newAccessToken(user)

	return LoginResult{
		AccessToken:  accessToken,
		TokenType:    "bearer",
		ExpiresIn:    exp,
		RefreshToken: newRefreshToken(user),
	}, nil
}

var authorizationRegexString = "^Bearer (?P<token>\\w+.\\w+.\\w+)$"
var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func DecodeJWT(authorization string) (*jwt.Token, error) {
	matches := authorizationRegex.FindStringSubmatch(authorization)
	if len(matches) != 2 {
		return nil, fmt.Errorf("unauthorized")
	}

	token, err := jwt.Parse(matches[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tokenHmacSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func newAccessToken(user entity.User) (string, int64) {
	claims := newBasicClaims(user, accessTokenExpiresIn)
	return newJwt(claims), claims["exp"].(int64)
}

func newRefreshToken(user entity.User) string {
	claims := newBasicClaims(user, refreshTokenExpiresIn)
	// TODO: add JWT ID (jti) for logout
	return newJwt(claims)
}

func newJwt(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenHmacSecret)
	if err != nil {
		return ""
	}
	return tokenString
}

func newBasicClaims(user entity.User, exp time.Duration) jwt.MapClaims {
	now := time.Now().UTC()
	return jwt.MapClaims{
		"iss": tokenIssuer,
		"nbf": now.Unix(),
		"exp": now.Add(exp).Unix(),

		"user_id":   user.ID,
		"user_name": user.Name,
	}
}
