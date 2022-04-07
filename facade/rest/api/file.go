package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	fileAggregate "github.com/yzx9/otodo/domain/file"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Upload public file, for user avatar
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "file required"))
		return
	}

	record, err := fileAggregate.UploadPublicFile(file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FileDTO{FileID: record.ID})
}

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	userID := common.MustGetAccessUserID(c)
	todoID, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "file required"))
		return
	}

	record, err := fileAggregate.UploadTodoFile(userID, todoID, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FileDTO{FileID: record.ID})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id := common.MustGetParam(c, "id")
	userID, err := common.GetAccessUserID(c)
	if err != nil {
		userID = 0
	}

	filepath, err := fileAggregate.GetFilePath(userID, id)
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
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "expiresIn required"))
		return
	}

	userID := common.MustGetAccessUserID(c)
	presigned, err := fileAggregate.CreateFilePreSignID(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.FilePreSignResultDTO{FileID: presigned})
}
