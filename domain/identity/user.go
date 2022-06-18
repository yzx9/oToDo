package identity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	id        int64
	createdAt time.Time
	updatedAt time.Time
	name      string
	nickname  string
	password  password
	email     string
	telephone string
	avatar    string
	githubId  int64

	fromSession string // login by session, fill with session id
}

func NewUser(
	id int64,
	createdAt time.Time,
	updatedAt time.Time,
	name string,
	nickname string,
	pwd password,
	email string,
	telephone string,
	avatar string,
	githubId int64,
) User {
	return User{
		id:        id,
		createdAt: createdAt,
		updatedAt: updatedAt,
		name:      name,
		nickname:  nickname,
		password:  pwd,
		email:     email,
		telephone: telephone,
		avatar:    avatar,
		githubId:  githubId,
	}
}

func (user User) Id() int64            { return user.id }
func (user User) CreatedAt() time.Time { return user.createdAt }
func (user User) UpdatedAt() time.Time { return user.updatedAt }
func (user User) Name() string         { return user.name }
func (user User) Nickname() string     { return user.nickname }
func (user User) Password() password   { return user.password }
func (user User) Email() string        { return user.email }
func (user User) Telephone() string    { return user.telephone }
func (user User) Avatar() string       { return user.avatar }
func (user User) GithubId() int64      { return user.githubId }

func (user *User) SetID(id int64) {
	if user.id == 0 {
		user.id = id
	}
}

func (user *User) new() (err error) {
	defer func() {
		if err != nil {
			err = Error{fmt.Errorf("fails to create user: %w", err)}
		}
	}()

	if err = UserRepository.Save(user); err != nil {
		return
	}

	PublishUserCreatedEvent(user.Id())

	return nil
}

func (user *User) Session() session {
	if user.fromSession == "" {
		// new session
		user.fromSession = uuid.NewString()
	}

	return session{
		userID:    user.Id(),
		sessionID: user.fromSession,
	}
}
