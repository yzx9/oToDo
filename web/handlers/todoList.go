package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Get todo lists
func GetTodoListsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodoLists(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get Todos by Todo List
func GetTodosFromTodoListHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithJson(c, "id is required")
		return
	}

	todos, err := bll.GetTodos(id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}
