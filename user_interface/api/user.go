package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/domain/aggregate/user"
	"github.com/yzx9/otodo/user_interface/common"
)

// Register
func PostUserHandler(c *gin.Context) {
	payload := dto.CreateUserDTO{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	user, err := user.CreateUser(payload)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}
