package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

func GetRequiredParam(c *gin.Context, name string) (string, error) {
	value, ok := c.Params.Get(name)
	if !ok {
		return "", util.NewError(errors.ErrPreconditionRequired, name+" required")
	}

	return value, nil
}

func GetRequiredParamID(c *gin.Context, name string) (int64, error) {
	value, err := GetRequiredParam(c, name)
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, util.NewError(errors.ErrPreconditionRequired, "invalid "+name)
	}

	return id, nil
}

func MustGetParam(c *gin.Context, name string) string {
	value, ok := c.Params.Get(name)
	if !ok {
		panic(name + " required")
	}

	return value
}
