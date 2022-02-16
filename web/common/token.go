package common

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

const AuthorizationHeaderKey = "Authorization"
const contextAccessTokenKey = "access_token"

func setAccessToken(c *gin.Context, token *jwt.Token) {
	c.Set(contextAccessTokenKey, token)
}

func GetAccessToken(c *gin.Context) (*jwt.Token, error) {
	authorization := c.Request.Header.Get(AuthorizationHeaderKey)
	token, err := bll.ParseAccessToken(authorization)
	if err != nil {
		return nil, util.NewError(otodo.ErrUnauthorized, "invalid token: %w", err)
	}

	setAccessToken(c, token)
	return token, nil
}

func GetAccessTokenClaims(c *gin.Context) (*dto.SessionTokenClaims, error) {
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.SessionTokenClaims)
	if !ok {
		return nil, util.NewError(otodo.ErrUnauthorized, "invalid token")
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
	value := c.MustGet(contextAccessTokenKey)
	token, _ := value.(*jwt.Token)
	return token
}

func MustGetAccessTokenClaims(c *gin.Context) *dto.SessionTokenClaims {
	token := MustGetAccessToken(c)
	claims, _ := token.Claims.(*dto.SessionTokenClaims)
	return claims
}

func MustGetAccessUserID(c *gin.Context) string {
	claims := MustGetAccessTokenClaims(c)
	return claims.UserID
}
