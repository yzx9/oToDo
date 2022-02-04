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

func ParseJWT(tokenString string, claims jwt.Claims) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tokenHmacSecret, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func NewJwt(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenHmacSecret)
	if err != nil {
		return ""
	}

	return tokenString
}

func NewClaims(exp time.Duration) jwt.StandardClaims {
	now := time.Now().UTC()
	return jwt.StandardClaims{
		Issuer:    tokenIssuer,
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		ExpiresAt: now.Add(exp).Unix(),
	}
}
