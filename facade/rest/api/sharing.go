package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/facade/rest/common"
)

// Get sharing info, only support todo list now
func GetSharingHandler(c *gin.Context) {
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharing, err := service.GetSharingInfo(token)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, sharing)
}

// Get todo list info by share token
func GetSharingTodoListHandler(c *gin.Context) {
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	info, err := service.GetSharingTodoListInfo(token)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, info)
}
