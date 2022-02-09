package bll

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Config
// TODO configurable
var tokenIssuer = "oToDo"
var tokenHmacSecret = []byte("test_secret")

type TokenClaims struct {
	jwt.StandardClaims
	UserID string `json:"uid"`
}

func NewClaims(userID string, exp time.Duration) TokenClaims {
	now := time.Now().UTC()
	return TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(exp).Unix(),
		},
		UserID: userID,
	}
}

func NewToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenHmacSecret)
	if err != nil {
		return ""
	}

	return tokenString
}

func ParseToken(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return tokenHmacSecret, nil
}
