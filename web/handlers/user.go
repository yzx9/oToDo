package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Register
func PostUserHandler(c *gin.Context) {
	payload := bll.CreateUserPayload{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	if len(payload.UserName) < 5 {
		c.String(http.StatusBadRequest, "invalid user name")
		return
	}

	if len(payload.Password) < 6 {
		c.String(http.StatusBadRequest, "password too short")
		return
	}

	user, err := bll.CreateUser(payload)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
