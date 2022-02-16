package dto

type SessionDTO struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type SessionTokenClaims struct {
	TokenClaims

	RefreshTokenID string `json:"rti,omitempty"`
	UserNickname   string `json:"nickname,omitempty"`
}

type LoginDTO struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}
