package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func HelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}
