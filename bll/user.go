package bll

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/entity"
)

func CreateInvalidUserRefreshToken(userID uuid.UUID, tokenID uuid.UUID) (entity.UserRefreshToken, error) {
	model, err := dal.InsertInvalidUserRefreshToken(entity.UserRefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenID:   tokenID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return entity.UserRefreshToken{}, fmt.Errorf("fails to invalid user refresh token, %w", err)
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
