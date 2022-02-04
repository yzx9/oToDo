package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/web/common"
)

func GetTodosHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithJson(c, "id required")
		return
	}

	todos, err := bll.GetTodo(id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

func PostTodoHandler(c *gin.Context) {
	todo := entity.Todo{}
	err := c.ShouldBind(&todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	todo, err = bll.AddTodo(todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func UpdateTodoHandler(c *gin.Context) {
	todo := entity.Todo{}
	err := c.ShouldBind(&todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	todo, err = bll.UpdateTodo(todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

func PatchTodoHandler(c *gin.Context) {
	// TODO How to do an patch
	// Can we get fields by reflect json object?
	c.Status(http.StatusNotImplemented)
}
