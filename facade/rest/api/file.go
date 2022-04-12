package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
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

	dto, err := service.UploadPublicFile(file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
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

	dto, err := service.UploadTodoFile(userID, todoID, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id := common.MustGetParam(c, "id")
	userID, err := common.GetAccessUserID(c)
	if err != nil {
		userID = 0
	}

	filepath, err := service.GetFilePath(userID, id)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.File(filepath)
}

// Create presigned file id for open access
func PostFilePreSignHandler(c *gin.Context) {
	payload := dto.FilePreSign{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "expiresIn required"))
		return
	}

	id, err := common.GetRequiredParamID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
		return
	}
	payload.FileID = id

	payload.UserID = common.MustGetAccessUserID(c)
	file, err := service.PreSignFile(payload)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, file)
}
