package handlers

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/utils"
	"github.com/yzx9/otodo/web/common"
)

type FilePayload struct {
	FileID string `json:"fileID"`
}

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

// Upload public file, for user avatar, only support img now
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "file required"))
		return
	}

	if !supportedFileTypeRegex.MatchString(file.Filename) {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionFailed, "invalid file type"))
		return
	}

	fileID, err := bll.UploadPublicFile(file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, FilePayload{fileID})
}

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	fileID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "file required"))
		return
	}

	fileID, err = bll.UploadTodoFile(fileID, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, FilePayload{fileID})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
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
	id, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	payload := struct {
		ExpiresIn int `json:"expiresIn"` // Unix
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, utils.NewError(otodo.ErrPreconditionRequired, "expiresIn required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	presigned, err := bll.CreateFilePresignedID(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, FilePayload{presigned})
}

// Get file from presigned id
func GetFilePresignedHandler(c *gin.Context) {
	presignedID, err := common.GetRequiredParam(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
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
