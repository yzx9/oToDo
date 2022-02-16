package entity

type UserInvalidRefreshToken struct {
	Entity

	UserID int64 `json:"userID"`
	User   User  `json:"-"`

	TokenID string `json:"tokenID" gorm:"type:char(36);"`
}
