package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/utils"
	webUtils "github.com/yzx9/otodo/web/utils"
)

// Ping Test
func GetSessionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}

// Login
func PostSessionHandler(c *gin.Context) {
	payload := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := bll.Login(payload.UserName, payload.Password)
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

// Logout, unactive refresh token if exists
func DeleteSessionHandler(c *gin.Context) {
	if token, err := parseRefreshToken(c); err == nil {
		bll.Logout(token)
	}

	c.String(http.StatusOK, "See you!")
}

// Refresh Current Access Token
func GetAccessTokenHandler(c *gin.Context) {
	userID, err := webUtils.GetAccessUserID(c)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	refreshAccessToken(c, userID)
}

// Create New Access Token by Refresh Token
func PostAccessTokenHandler(c *gin.Context) {
	token, err := parseRefreshToken(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	claims, ok := token.Claims.(*utils.TokenClaims)
	if !ok || !token.Valid {
		c.String(http.StatusBadRequest, "invalid token")
		return
	}

	refreshAccessToken(c, claims.UserID)
}

func parseRefreshToken(c *gin.Context) (*jwt.Token, error) {
	obj := &struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	if err := c.ShouldBind(&obj); err != nil {
		return nil, err
	}

	token, err := utils.ParseJWT(obj.RefreshToken)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func refreshAccessToken(c *gin.Context, userID string) {
	token, err := bll.NewAccessToken(userID)
	if err != nil {
		c.String(http.StatusBadRequest, "fails to refresh an token, %v", err.Error())
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}{token.AccessToken, token.TokenType, token.ExpiresIn})
}
