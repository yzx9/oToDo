package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/utils"
)

func ErrorMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Next()

		if c.IsAborted() {
			err := c.Errors.Last()
			typedError := utils.ErrorWithHttpStatus{}
			if errors.As(err, &typedError) {
				code := getHttpCodeFromError(typedError)
				c.AbortWithError(code, typedError)
				return
			}

			c.AbortWithError(http.StatusBadRequest, err)
		}
	}
}

func getHttpCodeFromError(err utils.ErrorWithHttpStatus) int {
	return err.Code
}
