package handler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
	"github.com/yzx9/otodo/web/common"
)

var supportedFileTypeRegex = regexp.MustCompile(`.(jpg|jpeg|JPG|png|PNG|gif|GIF|ico|ICO)$`)

// Upload public file, for user avatar, only support img now
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionRequired, "file required"))
		return
	}

	if !supportedFileTypeRegex.MatchString(file.Filename) {
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionFailed, "invalid file type"))
		return
	}

	fileID, err := bll.UploadPublicFile(file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FileDTO{FileID: fileID})
}

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	fileID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionRequired, "file required"))
		return
	}

	fileID, err = bll.UploadTodoFile(fileID, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FileDTO{FileID: fileID})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	userID, err := common.GetAccessUserID(c)
	if err != nil {
		// TODO check if this is an bug
		userID = 0
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
	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	payload := dto.FilePreSignDTO{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionRequired, "expiresIn required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	presigned, err := bll.CreateFilePresignedID(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FilePreSignResultDTO{FileID: presigned})
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
