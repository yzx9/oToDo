package dto

type UserCredential struct {
	UserName              string `json:"userName"`
	Password              string `json:"password"`
	RefreshTokenExpiresIn int    `json:"refreshTokenExp"`
}
