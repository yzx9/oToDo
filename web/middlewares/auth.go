package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

func JwtAuthMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		token, err := common.GetAccessToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid token")
			return
		}

		if bll.ShouldRefreshAccessToken(token) {
			claims := common.MustGetAccessTokenClaims(c)
			userID, _ := uuid.Parse(claims.UserID)
			refreshTokenID, _ := uuid.Parse(claims.RefreshTokenID)
			newToken, err := bll.NewAccessToken(userID, refreshTokenID)

			if err == nil {
				c.Header(common.AuthorizationHeaderKey, newToken.TokenType+" "+newToken.AccessToken)
			}
		}

		c.Next()
	}
}
