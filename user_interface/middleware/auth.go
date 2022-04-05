package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/domain/aggregate/user"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
	"github.com/yzx9/otodo/user_interface/common"
)

func JwtAuthMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := common.GetAccessToken(c)
		if err != nil {
			common.AbortWithError(c, util.NewError(errors.ErrUnauthorized, "invalid token"))
			return
		}

		if user.ShouldRefreshAccessToken(token) {
			claims := common.MustGetAccessTokenClaims(c)
			newToken, err := user.NewAccessToken(claims.UserID, claims.RefreshTokenID)

			if err == nil {
				c.Header(common.AuthorizationHeaderKey, newToken.TokenType+" "+newToken.AccessToken)
			}
		}

		c.Next()
	}
}
