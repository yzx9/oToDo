package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/api/common"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
)

// Get sharing info, only support todo list now
func GetSharingHandler(c *gin.Context) {
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharing, err := bll.ValidSharing(token)
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

	sharing, err := bll.ValidSharing(token)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	user, err := bll.GetUser(sharing.UserID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	list, err := bll.ForceGetTodoList(sharing.RelatedID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusNotImplemented, dto.SharingTodoList{
		UserNickname: user.Nickname,
		TodoListName: list.Name,
	})
}
