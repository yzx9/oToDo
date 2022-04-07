package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/domain/user"
	"github.com/yzx9/otodo/facade/rest/common"
)

// Register
func PostUserHandler(c *gin.Context) {
	payload := user.NewUser{}
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
