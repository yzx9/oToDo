package service

import (
	"fmt"

	"github.com/yzx9/otodo/domain/user"
	"github.com/yzx9/otodo/infrastructure/repository"
)

func GetUser(userID int64) (user.User, error) {
	u, err := repository.UserRepo.Find(userID)
	if err != nil {
		return user.User{}, fmt.Errorf("fails to get user: %w", err)
	}

	return u, nil
}

func CreateNewToken(refreshToken string) (user.SessionTokens, error) {
	token, err := user.ParseSessionToken(refreshToken)
	claims, ok := token.Claims.(*user.SessionTokenClaims)
	if err != nil || !ok || !token.Valid {
		return user.SessionTokens{}, fmt.Errorf("invalid token")
	}

	valid, err := user.IsValidRefreshToken(claims.UserID, claims.RefreshTokenID)
	if err != nil {
		return user.SessionTokens{}, err
	}

	if !valid {
		return user.SessionTokens{}, fmt.Errorf("refresh token has been invalid")
	}

	newToken, err := user.NewAccessToken(claims.UserID, claims.RefreshTokenID)
	if err != nil {
		return user.SessionTokens{}, fmt.Errorf("fails to refresh token: %w", err)
	}

	return newToken, nil
}
