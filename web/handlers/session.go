package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yzx9/otodo/bll"
	"github.com/yzx9/otodo/web/common"
)

// Ping Test
func GetSessionHandler(c *gin.Context) {
	claims := common.MustGetAccessTokenClaims(c)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("hello, %v", claims.UserNickname)})
}

// Login
func PostSessionHandler(c *gin.Context) {
	payload := struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
	}{}
	if err := c.ShouldBind(&payload); err != nil {
		common.AbortWithError(c, err)
		return
	}

	tokens, err := bll.Login(payload.UserName, payload.Password)
	if err != nil {
		common.AbortWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{tokens.AccessToken, tokens.TokenType, tokens.ExpiresIn, tokens.RefreshToken})
}

// Logout, unactive refresh token
func DeleteSessionHandler(c *gin.Context) {
	claims := common.MustGetAccessTokenClaims(c)
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		common.AbortWithError(c, fmt.Errorf("invalid user id: %v", claims.UserID))
	}

	refreshTokenID, err := uuid.Parse(claims.RefreshTokenID)
	if err != nil {
		common.AbortWithError(c, fmt.Errorf("invalid token id: %v", claims.RefreshTokenID))
	}

	err = bll.Logout(userID, refreshTokenID)
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

	if !bll.IsValidRefreshToken(userID, refreshTokenID) {
		common.AbortWithError(c, err)
		return
	}

	newToken, err := bll.NewAccessToken(userID, refreshTokenID)
	if err != nil {
		common.AbortWithJson(c, fmt.Sprintf("fails to refresh an token, %v", err.Error()))
		return
	}

	c.JSON(http.StatusOK, struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}{newToken.AccessToken, newToken.TokenType, newToken.ExpiresIn})
}

func parseRefreshToken(c *gin.Context) (uuid.UUID, uuid.UUID, error) {
	u := uuid.UUID{}
	obj := &struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	if err := c.ShouldBind(&obj); err != nil {
		return u, u, err
	}

	token, err := bll.ParseSessionToken(obj.RefreshToken)
	claims, ok := token.Claims.(*bll.SessionTokenClaims)
	if err != nil || !ok || !token.Valid {
		return u, u, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return u, u, fmt.Errorf("invalid user id: %v", claims.RefreshTokenID)
	}

	refreshTokenID, err := uuid.Parse(claims.RefreshTokenID)
	if err != nil {
		return u, u, fmt.Errorf("invalid refresh token id: %v", claims.RefreshTokenID)
	}

	return userID, refreshTokenID, nil
}
