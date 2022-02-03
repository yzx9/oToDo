package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/utils"
)

// Upload file, only support single file now
func PostFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "invalid file")
		return
	}

	filename, err := bll.UploadFile(file)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
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
		c.String(http.StatusBadRequest, "invalid file")
		return
	}

	userID := utils.MustGetAccessUserID(c)
	filepath, err := bll.GetFilePath(userID, params.id)
	if err != nil {
		return
	}

	c.File(filepath)
}
