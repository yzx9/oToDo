package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/application/dto"
	"github.com/yzx9/otodo/application/service"
	"github.com/yzx9/otodo/facade/rest/common"
	"github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/infrastructure/util"
)

// Ping Test
func GetSessionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "software is eating the world"})
}

// Login
func PostSessionHandler(c *gin.Context) {
	dto := dto.UserCredential{}
	if err := c.ShouldBind(&dto); err != nil {
		common.AbortWithError(c, err)
		return
	}

	tokens, err := service.Login(dto)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout, unactive refresh token
func DeleteSessionHandler(c *gin.Context) {
	token := common.MustGetAccessToken(c)
	service.Logout(token)
	c.JSON(http.StatusOK, gin.H{"message": "see you"})
}

// Create New Access Token by Refresh Token
func PostSessionTokenHandler(c *gin.Context) {
	payload := struct {
		RefreshToken string `json:"refreshToken"`
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, fmt.Errorf("refreshToken required"))
	}

	token, err := service.LoginByRefreshToken(payload.RefreshToken)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, token)
}

/**
 * OAuth
 */

func GetSessionOAuthGithub(c *gin.Context) {
	dto, err := service.CreateGithubOAuthURI()
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto)
}

func PostSessionOAuthGithub(c *gin.Context) {
	var payload dto.OAuthPayload
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, util.NewError(errors.ErrPreconditionRequired, "code, state required"))
		return
	}

	tokens, err := service.LoginByGithubOAuth(payload.Code, payload.State)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}
