package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

func GetTodoListsHandler(c *gin.Context) {
	userID := utils.MustGetAccessUserID(c)
	todos, err := bll.GetTodoLists(userID)
	if err != nil {
		utils.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodoListHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		utils.AbortWithJson(c, "id is required")
		return
	}

	todos, err := bll.GetTodos(id)
	if err != nil {
		utils.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}
