package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

func GetTodoListsHandler(c *gin.Context) {
	userID, err := utils.GetAccessUserID(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	todos, err := bll.GetTodoLists(userID)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, todos)
}

func GetTodoListHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.String(http.StatusPreconditionRequired, "id is required")
		return
	}

	todos, err := bll.GetTodos(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, todos)
}
