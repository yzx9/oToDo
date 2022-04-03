package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/user_interface/common"
	"github.com/yzx9/otodo/util"
)

func JwtAuthMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := common.GetAccessToken(c)
		if err != nil {
			common.AbortWithError(c, util.NewError(otodo.ErrUnauthorized, "invalid token"))
			return
		}

		if bll.ShouldRefreshAccessToken(token) {
			claims := common.MustGetAccessTokenClaims(c)
			newToken, err := bll.NewAccessToken(claims.UserID, claims.RefreshTokenID)

			if err == nil {
				c.Header(common.AuthorizationHeaderKey, newToken.TokenType+" "+newToken.AccessToken)
			}
		}

		c.Next()
	}
}
