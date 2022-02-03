package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/utils"
)

var accessTokenKey = "access_token"

func SetAccessToken(c *gin.Context, token *jwt.Token) {
	c.Set(accessTokenKey, token)
}

func GetAccessToken(c *gin.Context) (*jwt.Token, error) {
	value, exists := c.Get(accessTokenKey)
	if !exists {
		return nil, fmt.Errorf("access token not existed")
	}

	token, ok := value.(*jwt.Token)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func GetAccessTokenClaims(c *gin.Context) (*utils.TokenClaims, error) {
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*utils.TokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetAccessUserID(c *gin.Context) (string, error) {
	claims, err := GetAccessTokenClaims(c)
	if err != nil {
		return "", err
	}

	return claims.UserID, nil
}

func MustGetAccessToken(c *gin.Context) *jwt.Token {
	value := c.MustGet(accessTokenKey)
	token, _ := value.(*jwt.Token)
	return token
}

func MustGetAccessTokenClaims(c *gin.Context) *utils.TokenClaims {
	token := MustGetAccessToken(c)
	claims, _ := token.Claims.(*utils.TokenClaims)
	return claims
}

func MustGetAccessUserID(c *gin.Context) string {
	claims := MustGetAccessTokenClaims(c)
	return claims.UserID
}