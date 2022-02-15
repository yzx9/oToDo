package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/web/common"
)

// Create todo
func PostTodoHandler(c *gin.Context) {
	todo := entity.Todo{}
	if err := c.ShouldBind(&todo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err := bll.CreateTodo(userID, &todo); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Get todo
func GetTodoHandler(c *gin.Context) {
	todoID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.GetTodo(userID, todoID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}

// Update todo fully
func PutTodoHandler(c *gin.Context) {
	todo := entity.Todo{}
	err := c.ShouldBind(&todo)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err = bll.UpdateTodo(userID, &todo); err != nil {
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
	todoID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodo(userID, todoID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
