package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/bll"
)

func GetTodoListsHandler(c *gin.Context) {
	// TODO: get user by token
	todos, err := bll.GetTodoLists("")
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
