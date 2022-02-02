package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var tokenIssuer = "oToDo"
var tokenHmacSecret = []byte("test_secret") // TODO

type TokenClaims struct {
	jwt.StandardClaims

	UserID   string `json:"user_id"`
	UserName string `json:"user_name,omitempty"`
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return tokenHmacSecret, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, keyFunc)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func NewJwt(claims TokenClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(tokenHmacSecret)
	if err != nil {
		return ""
	}

	return tokenString
}

func NewTokenClaims(userID uuid.UUID, exp time.Duration) TokenClaims {
	now := time.Now().UTC()
	return TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    tokenIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(exp).Unix(),
		},
		UserID: userID.String(),
	}
}
