package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

// Get basic todo lists for current user
func GetCurrentUserBasicTodoListHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		TodoListID string `json:"todo_list_id"`
	}{user.BasicTodoListID.String()})
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

// Get Todos by Todo List
func GetTodosByTodoListHandler(c *gin.Context) {
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
