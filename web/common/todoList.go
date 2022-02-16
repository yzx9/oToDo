package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
)

func HandleGetCurrentUserTodoList(c *gin.Context, todoListID int64) {
	userID := MustGetAccessUserID(c)
	todoList, err := bll.SelectTodoList(userID, todoListID)
	if err != nil {
		AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, todoList)
}
