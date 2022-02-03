package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

func JwtAuthMiddleware(c *gin.Context) {
	const key = "Authorization"

	authorization := c.Request.Header.Get(key)
	token, err := bll.ParseAccessToken(authorization)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	utils.SetAccessToken(c, token)

	if bll.ShouldRefreshAccessToken(token) {
		userID := utils.MustGetAccessUserID(c)
		if newToken, err := bll.NewAccessToken(userID); err == nil {
			c.Header(key, newToken.TokenType+" "+newToken.AccessToken)
		}
	}

	c.Next()
}