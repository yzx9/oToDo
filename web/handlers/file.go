package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// TODO support create temporary opening aceess url

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	todoID, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
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
	}{fileID.String()})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	id, err := common.GetParamUUID(c, "id")
	if err != nil {
		common.AbortWithError(c, err)
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
