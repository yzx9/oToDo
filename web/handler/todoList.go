package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/web/common"
)

// Create todo list
func PostTodoListHandler(c *gin.Context) {
	list := entity.TodoList{}
	if err := c.ShouldBind(&list); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err := bll.CreateTodoList(userID, &list); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// Get todo list
func GetTodoListHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	common.HandleGetCurrentUserTodoList(c, id)
}

// Get todos in todo list
func GetTodoListTodosHandler(c *gin.Context) {
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodos(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Delete todo list
func DeleteTodoListHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodoList(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

/**
 * Todo List Sharing
 */

// Get shared users in todo list
func GetTodoListSharedUsersHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	users, err := bll.GetTodoListSharedUsers(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	// TODO[bug]: filter user fields

	c.JSON(http.StatusOK, users)
}

// Delete shared user from todo list
func DeleteTodoListSharedUserHandler(c *gin.Context) {
	operatorID := common.MustGetAccessUserID(c)
	todoListID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID, err := common.GetRequiredParamID(c, "user-id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if err := bll.DeleteTodoListSharedUser(operatorID, userID, todoListID); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// Create share link for todo list
func PostTodoListSharingsHandler(c *gin.Context) {
	// TODO 创建分享链接，会使之前的失效
	c.Status(http.StatusNotImplemented)
}

// Get current share link
func GetTodoListSharingsHandler(c *gin.Context) {
	// TODO 获取分享连接
	c.Status(http.StatusNotImplemented)
}

// Join todo list by share link
func PostTodoListSharingHandler(c *gin.Context) {
	// TODO 加入共享列表，传入token
	c.Status(http.StatusNotImplemented)
}

// Inactive share link
func DeleteTodoListSharingHandler(c *gin.Context) {
	// TODO 删除分享链接
	c.Status(http.StatusNotImplemented)
}
