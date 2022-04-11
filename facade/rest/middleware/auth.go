package middleware

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

const ContextKeyAccessToken = "Authorization"
const ContextKeyUserID = "UserID"
const headerKey = "Authorization"

var regex = regexp.MustCompile(`^[Bb]earer (?P<token>[\w-]+.[\w-]+.[\w-]+)$`)

func JwtAuthMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get(headerKey)
		matches := regex.FindStringSubmatch(authorization)
		if len(matches) != 2 {
			c.Error(util.NewError(errors.ErrUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		token := matches[1]
		validation := service.LoginByAccessToken(token)
		if !validation.Valid {
			c.Error(util.NewError(errors.ErrUnauthorized, "invalid token"))
			c.Abort()
			return
		}

		c.Keys[ContextKeyAccessToken] = token
		c.Keys[ContextKeyUserID] = validation.UserID
		if validation.NewAccessToken {
			c.Header(headerKey, "bearer "+validation.AccessToken)
		}

		c.Next()
	}
}
