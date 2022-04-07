package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/facade/rest/common"
)

// Get sharing info, only support todo list now
func GetSharingHandler(c *gin.Context) {
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharing, err := todolist.ValidSharing(token)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.SharingToken{
		Token:     sharing.Token,
		Type:      sharing.Type,
		CreatedAt: sharing.CreatedAt,
	})
}

// Get todo list info by share token
func GetSharingTodoListHandler(c *gin.Context) {
	token, err := common.GetRequiredParam(c, "token")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	sharing, err := todolist.ValidSharing(token)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	user, err := service.GetUser(sharing.UserID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	list, err := service.ForceGetTodoList(sharing.RelatedID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusNotImplemented, dto.SharingTodoList{
		UserNickname: user.Nickname,
		TodoListName: list.Name,
	})
}
