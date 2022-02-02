package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/bll"
)

func GetTodosHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		c.String(http.StatusPreconditionRequired, "id required")
		return
	}

	todos, err := bll.GetTodo(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, todos)
}

func CreateTodosHandler(c *gin.Context) {
	// TODO
	c.String(http.StatusNotImplemented, "todo")
}
