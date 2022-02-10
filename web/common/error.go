package common

import (
	"github.com/gin-gonic/gin"
)

func AbortWithError(c *gin.Context, err error) *gin.Error {
	gError := c.Error(err)
	c.Abort()
	return gError
}
