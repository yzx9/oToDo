package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/todo"
	todoAggregate "github.com/yzx9/otodo/domain/todo"
	"github.com/yzx9/otodo/facade/rest/common"
)

// Create todo
func PostTodoHandler(c *gin.Context) {
	todo := todo.Todo{}
	if err := c.ShouldBind(&todo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err := todoAggregate.CreateTodo(userID, &todo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Get todo
func GetTodoHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := service.GetTodo(userID, todoID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Update todo fully
func PutTodoHandler(c *gin.Context) {
	todo := todo.Todo{}
	err := c.ShouldBind(&todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err = todoAggregate.UpdateTodo(userID, &todo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Update todo partial
func PatchTodoHandler(c *gin.Context) {
	// TODO How to do an patch
	// Can we get fields by reflect json object?
	c.Status(http.StatusNotImplemented)
}

// Delete Todo
func DeleteTodoHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := todoAggregate.DeleteTodo(userID, todoID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
