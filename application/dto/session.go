package dto

type SessionTokens struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type SessionValidation struct {
	Valid          bool   `json:"valid"`
	UserID         int64  `json:"userID"`
	NewAccessToken bool   `json:"newAccessToken"`
	AccessToken    string `json:"token"`
}
