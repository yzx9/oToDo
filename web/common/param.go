package common

import (
	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
)

func GetRequiredParam(c *gin.Context, name string) (string, error) {
	value, ok := c.Params.Get(name)
	if !ok {
		return "", utils.NewError(otodo.ErrPreconditionRequired, name+" required")
	}

	return value, nil
}
