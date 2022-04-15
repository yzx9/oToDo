package session

import "github.com/yzx9/otodo/domain/user"

var UserRepository userRepository
var UserInvalidRefreshTokenRepository userInvalidRefreshTokenRepository

type userRepository interface {
	Find(id int64) (user.User, error)

	FindByUserName(username string) (user.User, error)

	FindByGithubID(githubID int64) (user.User, error)
}

type userInvalidRefreshTokenRepository interface {
	Save(entity *UserInvalidRefreshToken) error

	Exist(userID int64, tokenID string) (bool, error)
}
