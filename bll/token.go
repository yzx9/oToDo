package bll

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/infrastructure/config"
	"github.com/yzx9/otodo/model/dto"
)

func NewClaims(userID int64, exp time.Duration) dto.TokenClaims {
	now := time.Now().UTC()
	return dto.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    config.Secret.TokenIssuer,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(exp).Unix(),
		},
		UserID: userID,
	}
}

func NewToken(claims jwt.Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.Secret.TokenHmacSecret)
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
	return config.Secret.TokenHmacSecret, nil
}
