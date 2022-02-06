package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func CreateInvalidUserRefreshToken(userID string, tokenID string) (entity.UserRefreshToken, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("invalid user id: %v", userID)
	}

	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("invalid token id: %v", tokenID)
	}

	model := entity.UserRefreshToken{
		ID:        uuid.New(),
		UserID:    userUUID,
		TokenID:   tokenUUID,
		CreatedAt: time.Now(),
	}
	model, err = dal.InsertInvalidUserRefreshToken(model)
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("fails to insert invalid user refresh token")
	}

	return model, nil
}

// Verify is it an valid token.
// Note: This func don't check token expire time
func IsValidRefreshToken(userID string, tokenID string) bool {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return false
	}

	tokenUUID, err := uuid.Parse(tokenID)
	if err != nil {
		return false
	}

	return !dal.ExistInvalidUserRefreshToken(userUUID, tokenUUID)
}
