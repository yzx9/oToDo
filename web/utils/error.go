package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/utils"
)

func AbortWithError(c *gin.Context, err error) {
	if (errors.Is(err, &utils.ErrorWithHttpStatus{})) {
		return
	}

	c.AbortWithError(http.StatusBadRequest, err)
}

func AbortWithJson(c *gin.Context, jsonObj interface{}) {
	c.AbortWithStatusJSON(http.StatusBadRequest, jsonObj)
}
