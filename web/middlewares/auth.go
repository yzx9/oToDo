package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

func JwtAuthMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	token, err := bll.ParseAccessToken(authorization)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	utils.SetAccessToken(c, token)
	c.Next()
}
