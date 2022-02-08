package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Get basic todo lists for current user
func GetCurrentUserHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Get todo lists for current user
func GetCurrentUserTodoListsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todos, err := bll.GetTodoLists(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get todo list
func GetTodoListHandler(c *gin.Context) {
	id, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	getTodoListHandler(c, id)
}

// Get basic todo lists for current user
func GetCurrentUserBasicTodoListHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	getTodoListHandler(c, user.BasicTodoListID)
}

func getTodoListHandler(c *gin.Context, todoListID uuid.UUID) {
	userID := common.MustGetAccessUserID(c)
	todoList, err := bll.GetTodoList(userID, todoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todoList)
}

// Get todos in todo list
func GetTodoListTodosHandler(c *gin.Context) {
	id, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	todos, err := bll.GetTodos(id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}
