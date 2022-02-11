package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
	"github.com/yzx9/otodo/web/common"
)

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	todoID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, fmt.Errorf("id required"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, err)
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
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	filepath, err := bll.GetFilePath(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.File(filepath)
}

// Create presigned file id for open access
func PostFilePresignHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "id required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	presigned := bll.CreateFilePresignedID(userID, id)
	c.JSON(http.StatusOK, struct {
		PresignedFileID string `json:"presigned_file_id"`
	}{presigned})
}

// Get file from presigned id
func GetFilePresignedHandler(c *gin.Context) {
	presignedID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "id required"))
		return
	}

	id, err := bll.ParseFileSignedID(presignedID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	filepath, err := bll.ForceGetFilePath(id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.File(filepath)
}
