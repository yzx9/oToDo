package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/utils"
	"github.com/yzx9/otodo/web/common"
)

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	todoID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("id required"))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("file required"))
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
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("id required"))
		return
	}

	userID, err := common.GetAccessUserID(c)
	if err != nil {
		userID = ""
	}

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
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("id required"))
		return
	}

	payload := struct {
		ExpiresIn int `json:"expires_in"` // Unix
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("expires_in required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	presigned, err := bll.CreateFilePresignedID(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		PresignedFileID string `json:"presigned_file_id"`
	}{presigned})
}

// Get file from presigned id
func GetFilePresignedHandler(c *gin.Context) {
	presignedID, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithError(c, utils.NewErrorWithPreconditionRequired("id required"))
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
