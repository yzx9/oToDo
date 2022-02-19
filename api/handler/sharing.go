package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Get sharing info, only support todo list now
func GetSharingHandler(c *gin.Context) {
	// TODO 获取共享链接信息
	c.Status(http.StatusNotImplemented)
}
