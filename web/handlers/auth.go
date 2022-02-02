package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
)

func LoginHandler(c *gin.Context) {
	// TODO
	userName := "admin"
	password := "admin123"

	tokens, err := bll.Login(userName, password)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{tokens.AccessToken, tokens.TokenType, tokens.ExpiresIn, tokens.RefreshToken})
}

func LogoutHandler(c *gin.Context) {
	c.String(http.StatusOK, "See you!")
}
