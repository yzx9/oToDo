package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/driver/rest/common"
)

// Create todo
func PostTodoHandler(c *gin.Context) {
	newTodo := dto.NewTodo{}
	if err := c.ShouldBind(&newTodo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := service.CreateTodo(userID, newTodo)
	if err != nil {
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
	todo := dto.Todo{}
	err := c.ShouldBind(&todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	newTodo, err := service.UpdateTodo(userID, todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, newTodo)
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
	todo, err := service.DeleteTodo(userID, todoID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
