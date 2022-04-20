package dto

import "github.com/yzx9/otodo/domain/identity"

type NewUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}

func (u NewUser) ToEntity() identity.NewUser {
	return identity.NewUser{
		UserName: u.UserName,
		Password: u.Password,
		Nickname: u.Nickname,
	}
}
