package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/entity"
	"github.com/yzx9/otodo/web/common"
)

// Get current user
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
	todos, err := bll.SelectTodoLists(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get basic todo lists for current user
func GetCurrentUserBasicTodoListHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := bll.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	common.HandleGetCurrentUserTodoList(c, user.BasicTodoListID)
}

// Get daily todos for current user
func GetCurrentUserDailyTodosHandler(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// Get planned todos for current user
func GetCurrentUserPlannedTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, bll.GetPlannedTodos)
}

// Get important todos for current user
func GetCurrentUserImportantTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, bll.GetImportantTodos)
}

// Get not-notified todos for current user
func GetCurrentUserNotNotifiedTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, bll.GetNotNotifiedTodos)
}

func handleGetCurrentUserTodos(c *gin.Context, getTodos func(string) ([]entity.Todo, error)) {
	userID := common.MustGetAccessUserID(c)
	todos, err := getTodos(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get todo list folders for current user
func GetCurrentUserTodoListFoldersHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	folders, err := bll.GetTodoListFolders(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folders)
}
