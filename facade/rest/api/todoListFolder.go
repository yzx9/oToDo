package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/domain/todolist"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/repository"
)

// Create todo list folder
func PostTodoListFolderHandler(c *gin.Context) {
	folder := repository.TodoListFolder{}
	if err := c.ShouldBind(&folder); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	if err := todolist.CreateTodoListFolder(userID, &folder); err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Get todo list folder
func GetTodoListFolderHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	folder, err := service.GetTodoListFolder(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Delete todo list folder
func DeleteTodoListFolderHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := todolist.DeleteTodoListFolder(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
