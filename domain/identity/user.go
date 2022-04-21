package identity

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Nickname  string
	Password  []byte
	Email     string
	Telephone string
	Avatar    string
	GithubID  int64
}

type NewUser struct {
	UserName string
	Password string
	Nickname string
}

func CreateUser(payload NewUser) (User, error) {
	if len(payload.UserName) < 5 {
		return User{}, UserNameTooShort
	}

	if len(payload.Password) < 6 {
		return User{}, PasswordTooShort
	}

	exist, err := UserRepository.ExistByUserName(payload.UserName)
	if err != nil {
		return User{}, newErr(fmt.Errorf("fails to valid user name: %w", err))
	}

	if exist {
		return User{}, UserNameDuplicate
	}

	user := User{
		Name:     payload.UserName,
		Nickname: payload.Nickname,
	}
	user.Password = user.cryptoPassword(payload.Password)
	if err := user.new(); err != nil {
		return User{}, newErr(fmt.Errorf("fails to create user: %w", err))
	}

	return user, nil
}

// Password
func (user User) ValidatePassword(password string) bool {
	crypto := user.cryptoPassword(password)
	return bytes.Equal(user.Password, crypto)
}

func (User) cryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), Conf.PasswordNonce...))
	return pwd[:]
}

/**
 * OAuth
 */

func GetOrRegisterUserByGithub(profile GithubUserPublicProfile) (User, error) {
	exist, err := UserRepository.ExistByGithubID(profile.ID)
	if err != nil {
		return User{}, newErr(fmt.Errorf("fails to register user: %w", err))
	}

	if exist {
		user, err := UserRepository.FindByGithubID(profile.ID)
		if err != nil {
			return User{}, newErr(fmt.Errorf("fails to get user: %w", err))
		}

		return user, nil
	}

	// Register new user
	// TODO[feat]: download user avatar
	user := User{
		Name:     profile.Email,
		Nickname: profile.Name,
		Email:    profile.Email,
		GithubID: profile.ID,
	}
	if err := user.new(); err != nil {
		return User{}, newErr(fmt.Errorf("fails to get user: %w", err))
	}

	return user, nil
}

/**
 * Helpers
 */

func (user *User) new() (err error) {
	defer func() {
		if err != nil {
			err = newErr(fmt.Errorf("fails to create user: %w", err))
		}
	}()

	if err = UserRepository.Save(user); err != nil {
		return
	}

	PublishUserCreatedEvent(user.ID)

	return nil
}
