package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yzx9/otodo/bll"
)

func GetTodosHandler(c *gin.Context) {
	c.JSON(http.StatusOK, bll.GetTodos())
}
