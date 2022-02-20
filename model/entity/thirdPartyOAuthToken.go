package entity

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

type ThirdPartyOAuthToken struct {
	Entity

	Active bool   `json:"active"`
	Type   int8   `json:"type" gorm:"index:idx_third_party_oauth_tokens_user,unique"`
	Token  string `json:"token" gorm:"size:128"`
	Scope  string `json:"scope" gorm:"size:32"`

	UserID int64 `json:"userID" gorm:"index:idx_third_party_oauth_tokens_user,unique"`
	User   User  `json:"-"`
}

func (ThirdPartyOAuthToken) TableName() string {
	return "third_party_oauth_tokens"
}
