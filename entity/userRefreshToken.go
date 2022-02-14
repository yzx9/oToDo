package entity

type UserInvalidRefreshToken struct {
	Entity

	UserID  string `json:"user_id" gorm:"size:36"`
	TokenID string `json:"token_id" gorm:"size:36"`
}