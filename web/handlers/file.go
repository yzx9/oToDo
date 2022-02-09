package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// TODO support create temporary opening aceess url

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	todoID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithJson(c, "invalid file")
		return
	}

	fileID, err := bll.UploadTodoFile(todoID, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		FileID string `json:"file_id"`
	}{fileID})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	filepath, err := bll.GetFilePath(userID, id)
	if err != nil {
		common.AbortWithJson(c, "invalid file")
		return
	}

	c.File(filepath)
}
