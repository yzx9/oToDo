package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/otodo"
	"github.com/yzx9/otodo/util"
	"github.com/yzx9/otodo/web/common"
)

// Upload public file, for user avatar
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, util.NewError(otodo.ErrPreconditionRequired, "file required"))
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
	id := common.MustGetParam(c, "id")
	userID, err := common.GetAccessUserID(c)
	if err != nil {
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
func PostFilePreSignHandler(c *gin.Context) {
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
	presigned, err := bll.CreateFilePreSignID(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FilePreSignResultDTO{FileID: presigned})
}
