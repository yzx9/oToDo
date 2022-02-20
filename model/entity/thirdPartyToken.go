package entity

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

type ThirdPartyToken struct {
	Entity

	Active    bool   `json:"active"`
	Type      int8   `json:"type"`
	Token     string `json:"token" gorm:"size:128"`
	TokenType string `json:"tokenType" gorm:"size:16"`
	Scope     string `json:"scope" gorm:"size:32"`

	UserID int64 `json:"userID"`
	User   User  `json:"-"`
}
