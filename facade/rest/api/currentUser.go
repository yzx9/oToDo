package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/repository"
)

// Get current user
func GetCurrentUserHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := service.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Get menu
func GetCurrentUserMenu(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	menu, err := service.GetMenu(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, menu)
}

// Get todo lists for current user
func GetCurrentUserTodoListsHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todos, err := service.GetTodoLists(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get basic todo list todos for current user
func GetCurrentUserBasicTodoListTodosHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	user, err := service.GetUser(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	todos, err := service.ForceGetTodos(user.BasicTodoListID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todos)
}

// Get daily todos for current user
func GetCurrentUserDailyTodosHandler(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}

// Get planned todos for current user
func GetCurrentUserPlannedTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, service.GetPlannedTodos)
}

// Get important todos for current user
func GetCurrentUserImportantTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, service.GetImportantTodos)
}

// Get not-notified todos for current user
func GetCurrentUserNotNotifiedTodosHandler(c *gin.Context) {
	handleGetCurrentUserTodos(c, service.GetNotNotifiedTodos)
}

func handleGetCurrentUserTodos(c *gin.Context, getTodos func(userID int64) ([]repository.Todo, error)) {
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
	folders, err := service.GetTodoListFolders(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folders)
}
