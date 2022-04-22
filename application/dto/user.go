package dto

type NewUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
}
