package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Get todo lists folder for current user
func GetCurrentUserTodoListFoldersHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	folders, err := bll.GetTodoListFolders(userID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folders)
}

// Get todo list folder
func GetTodoListFolderHandler(c *gin.Context) {
	todoListFolderID, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	folder, err := bll.GetTodoListFolder(userID, todoListFolderID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, folder)
}

// Delete todo list folder
func DeleteTodoListFolderHandler(c *gin.Context) {
	id, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	todo, err := bll.DeleteTodoList(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
