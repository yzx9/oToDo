package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
)

func JwtAuthMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	token, err := bll.DecodeJWT(authorization)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	c.Set("access_token", token)
	c.Next()
}
