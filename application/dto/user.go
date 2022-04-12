package dto

import "github.com/yzx9/otodo/domain/user"

type NewUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func (u NewUser) ToEntity() user.NewUser {
	return user.NewUser{
		UserName: u.UserName,
		Password: u.Password,
		Nickname: u.Nickname,
	}
}
