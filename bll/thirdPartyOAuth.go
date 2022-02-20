package bll

import (
	"fmt"

	"github.com/yzx9/otodo/dal"
	"github.com/yzx9/otodo/model/entity"
)

func UpdateThirdPartyOAuthToken(token *entity.ThirdPartyOAuthToken) error {
	// TODO[bug]: handle error
	exist, err := dal.ExistActiveThirdPartyOAuthToken(token.UserID, entity.ThirdPartyTokenType(token.Type))
	if err != nil {
		return fmt.Errorf("fails to update third party oauth token: %w", err)
	}

	handle := dal.UpdateThirdPartyOAuthToken
	if !exist {
		handle = dal.InsertThirdPartyOAuthToken
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
