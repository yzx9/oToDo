package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

func GetTodosHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		utils.AbortWithJson(c, "id required")
		return
	}

	todos, err := bll.GetTodo(id)
	if err != nil {
		utils.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func CreateTodosHandler(c *gin.Context) {
	// TODO
	c.String(http.StatusNotImplemented, "todo")
}
