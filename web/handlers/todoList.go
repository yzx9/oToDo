package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/entity"
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
	id, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	common.HandleGetCurrentUserTodoList(c, id)
}

// Get todos in todo list
func GetTodoListTodosHandler(c *gin.Context) {
	todoListID, err := common.GetRequiredParam(c, "id")
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
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
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
