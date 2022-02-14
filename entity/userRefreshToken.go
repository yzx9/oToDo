package entity

type UserInvalidRefreshToken struct {
	Entity

	UserID string `json:"user_id" gorm:"size:36"`
	User   User   `json:"-"`

	TokenID string `json:"token_id" gorm:"size:36"`
}
