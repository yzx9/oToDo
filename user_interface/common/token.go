package common

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

const AuthorizationHeaderKey = "Authorization"
const contextAccessTokenKey = "access_token"

func GetAccessToken(c *gin.Context) (*jwt.Token, error) {
	authorization := c.Request.Header.Get(AuthorizationHeaderKey)
	token, err := bll.ParseAccessToken(authorization)
	if err != nil {
		return nil, util.NewError(errors.ErrUnauthorized, "invalid token: %w", err)
	}

	c.Set(contextAccessTokenKey, token)
	return token, nil
}

func GetAccessTokenClaims(c *gin.Context) (*dto.SessionTokenClaims, error) {
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.SessionTokenClaims)
	if !ok {
		return nil, util.NewError(errors.ErrUnauthorized, "invalid token")
	}

	return claims, nil
}

func GetAccessUserID(c *gin.Context) (int64, error) {
	claims, err := GetAccessTokenClaims(c)
	if err != nil {
		return 0, err
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

func MustGetAccessUserID(c *gin.Context) int64 {
	claims := MustGetAccessTokenClaims(c)
	return claims.UserID
}
