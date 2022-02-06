package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Upload todo file, only support single file now
func PostTodoFileHandler(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		common.AbortWithJson(c, "id required")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		common.AbortWithJson(c, "invalid file")
		return
	}

	filename, err := bll.UploadTodoFile(id, file)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		FileName string `json:"file_name"`
	}{filename})
}

// Upload file, only support single file now
func GetFileHandler(c *gin.Context) {
	params := struct{ id string }{}
	err := c.ShouldBind(&params)
	if err != nil {
		common.AbortWithJson(c, "invalid file")
		return
	}

	userID := common.MustGetAccessUserID(c)
	filepath, err := bll.GetFilePathWithAuth(params.id, userID)
	if err != nil {
		common.AbortWithJson(c, "invalid file")
		return
	}

	c.File(filepath)
}
