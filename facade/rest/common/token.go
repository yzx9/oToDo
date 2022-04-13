package common

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/facade/rest/middleware"
)

const AuthorizationHeaderKey = "Authorization"
const authorizationRegexString = `^[Bb]earer (?P<token>[\w-]+.[\w-]+.[\w-]+)$`

var authorizationRegex = regexp.MustCompile(authorizationRegexString)

func GetAccessToken(c *gin.Context) (string, error) {
	token, ok := c.Get(middleware.ContextKeyAccessToken)
	if !ok {
		return "", fmt.Errorf("unauthorized")
	}

	return token.(string), nil
}

func GetAccessUserID(c *gin.Context) (int64, error) {
	userID, ok := c.Get(middleware.ContextKeyUserID)
	if !ok {
		return 0, fmt.Errorf("unauthorized")
	}

	return userID.(int64), nil
}

func MustGetAccessToken(c *gin.Context) string {
	token, ok := c.Get(middleware.ContextKeyAccessToken)
	if !ok {
		return ""
	}

	return token.(string)
}

func MustGetAccessUserID(c *gin.Context) int64 {
	userID, ok := c.Get(middleware.ContextKeyUserID)
	if !ok {
		return 0
	}

	return userID.(int64)
}
