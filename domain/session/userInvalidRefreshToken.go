package session

import (
	"fmt"
	"time"
)

type UserInvalidRefreshToken struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID  int64
	TokenID string
}

func CreateUserInvalidRefreshToken(userID int64, tokenID string) (UserInvalidRefreshToken, error) {
	model := UserInvalidRefreshToken{
		UserID:  userID,
		TokenID: tokenID,
	}
	if err := UserInvalidRefreshTokenRepository.Save(&model); err != nil {
		return UserInvalidRefreshToken{}, fmt.Errorf("fails to make user refresh token invalid: %w", err)
	}

	return model, nil
}
