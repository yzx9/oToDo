package common

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/bll"
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
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	setAccessToken(c, token)
	return token, nil
}

func GetAccessTokenClaims(c *gin.Context) (*bll.SessionTokenClaims, error) {
	token, err := GetAccessToken(c)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*bll.SessionTokenClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetAccessUserID(c *gin.Context) (uuid.UUID, error) {
	claims, err := GetAccessTokenClaims(c)
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user id in access token")
	}

	return id, nil
}

func MustGetAccessToken(c *gin.Context) *jwt.Token {
	value := c.MustGet(contextAccessTokenKey)
	token, _ := value.(*jwt.Token)
	return token
}

func MustGetAccessTokenClaims(c *gin.Context) *bll.SessionTokenClaims {
	token := MustGetAccessToken(c)
	claims, _ := token.Claims.(*bll.SessionTokenClaims)
	return claims
}

func MustGetAccessUserID(c *gin.Context) string {
	claims := MustGetAccessTokenClaims(c)
	return claims.UserID
}
