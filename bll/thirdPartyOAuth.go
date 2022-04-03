package bll

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/repository"
	"github.com/yzx9/otodo/model/entity"
)

func UpdateThirdPartyOAuthToken(token *entity.ThirdPartyOAuthToken) error {
	// TODO[bug]: handle error
	exist, err := repository.ExistActiveThirdPartyOAuthToken(token.UserID, entity.ThirdPartyTokenType(token.Type))
	if err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	handle := repository.UpdateThirdPartyOAuthToken
	if !exist {
		handle = repository.InsertThirdPartyOAuthToken
	}

	if err := handle(token); err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	return nil
}

func UpdateThirdPartyOAuthTokenAsync(token *entity.ThirdPartyOAuthToken) {
	if err := UpdateThirdPartyOAuthToken(token); err != nil {
		// TODO[bug]: handle error
		fmt.Println(err)
	}
}
