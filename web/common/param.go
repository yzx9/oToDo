package common

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
)

func MustGetParam(c *gin.Context, name string) string {
	value, ok := c.Params.Get(name)
	if !ok {
		panic(name + " required")
	}
	return value
}

func GetRequiredParam(c *gin.Context, name string) (string, error) {
	value, ok := c.Params.Get(name)
	if !ok {
		return "", util.NewError(otodo.ErrPreconditionRequired, name+" required")
	}

	return value, nil
}
