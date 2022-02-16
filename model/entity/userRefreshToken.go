package entity

type UserInvalidRefreshToken struct {
	Entity

	UserID string `json:"userID" gorm:"type:char(36);"`
	User   User   `json:"-"`

	TokenID string `json:"tokenID" gorm:"type:char(36);"`
}
