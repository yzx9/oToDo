package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/driver/rest/common"
)

// Create todo list folder
func PostTodoListFolderHandler(c *gin.Context) {
	newFolder := dto.NewTodoListFolder{}
	if err := c.ShouldBind(&newFolder); err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID := common.MustGetAccessUserID(c)
	folder, err := service.CreateTodoListFolder(userID, newFolder)
	if err != nil {
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
	todo, err := service.DeleteTodoListFolder(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todo)
}
