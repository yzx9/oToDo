package identity

import "fmt"

func GetUserByUserName(userName string) (User, error) {
	user, err := UserRepository.FindByUserName(userName)
	if err != nil {
		return User{}, InvalidCredential
	}

	return user, nil
}

func GetUserByGithubOAuth(code, state string) (User, error) {
	oauth, err := getOAuthEntryByState(state)
	if err != nil {
		return User{}, InvalidCredential
	}

	user, err := oauth.GetOrRegisterUserByGithub(code)
	if err != nil {
		return User{}, Error{fmt.Errorf("fails to get user: %w", err)}
	}

	return user, nil
}

func GetUserBySessionToken(tokenType sessionTokenType, token string) (User, error) {
	claims, err := parseSessionToken(token)
	if err != nil {
		return User{}, InvalidCredential
	}

	if tokenType == RefreshToken {
		valid, err := UserInvalidRefreshTokenRepository.Exist(claims.UserID, claims.Id)
		if err != nil || valid {
			return User{}, InvalidCredential
		}
	}

	user, err := UserRepository.Find(claims.UserID)
	if err != nil {
		return User{}, InvalidCredential
	}

	user.fromSession = claims.SessionID
	return user, nil
}

func RegisterUser(userName, nickname, pwd string) (User, error) {
	if len(userName) < 5 {
		return User{}, UserNameTooShort
	}

	password, err := NewPassword(pwd)
	if err != nil {
		return User{}, err
	}

	if password.Empty() {
		return User{}, InvalidPassword
	}

	exist, err := UserRepository.ExistByUserName(userName)
	if err != nil {
		return User{}, Error{fmt.Errorf("fails to valid user name: %w", err)}
	}

	if exist {
		return User{}, UserNameDuplicate
	}

	user := User{
		name:     userName,
		nickname: nickname,
		password: password,
	}
	if err := user.new(); err != nil {
		return User{}, Error{fmt.Errorf("fails to create user: %w", err)}
	}

	return user, nil
}

func NewGithubOAuthURI() (string, error) {
	oauth, err := newOAuthEntry()
	if err != nil {
		return "", Error{fmt.Errorf("fails to create github oauth uri: %w", err)}
	}

	uri, err := GithubAdapter.CreateOAuthURI(oauth.state)
	if err != nil {
		return "", Error{fmt.Errorf("fails to create github oauth uri: %w", err)}
	}

	return uri, nil
}
