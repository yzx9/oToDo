package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/web/common"
)

// Ping Test
func GetSessionHandler(c *gin.Context) {
	claims := common.MustGetAccessTokenClaims(c)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("hello, %v", claims.UserNickname)})
}

// Login
func PostSessionHandler(c *gin.Context) {
	payload := dto.LoginDTO{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	tokens, err := bll.Login(payload.UserName, payload.Password)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout, unactive refresh token
func DeleteSessionHandler(c *gin.Context) {
	claims := common.MustGetAccessTokenClaims(c)
	err := bll.Logout(claims.UserID, claims.RefreshTokenID)
	if err != nil {
		// TODO log
		fmt.Println(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "see you"})
}

// Create New Access Token by Refresh Token
func PostSessionTokenHandler(c *gin.Context) {
	userID, refreshTokenID, err := parseRefreshToken(c)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	valid, err := bll.IsValidRefreshToken(userID, refreshTokenID)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	if !valid {
		common.AbortWithError(c, fmt.Errorf("refresh token has been invalid"))
		return
	}

	newToken, err := bll.NewAccessToken(userID, refreshTokenID)
	if err != nil {
		msg := fmt.Sprintf("fails to refresh an token, %v", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, msg)
		return
	}

	c.JSON(http.StatusOK, newToken)
}

func parseRefreshToken(c *gin.Context) (int64, string, error) {
	obj := dto.RefreshTokenDTO{}
	if err := c.ShouldBind(&obj); err != nil {
		return 0, "", fmt.Errorf("refresh_token required")
	}

	token, err := bll.ParseSessionToken(obj.RefreshToken)
	claims, ok := token.Claims.(*dto.SessionTokenClaims)
	if err != nil || !ok || !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}

	return claims.UserID, claims.RefreshTokenID, nil
}
