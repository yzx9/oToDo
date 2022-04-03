package dto

type SessionToken struct {
	AccessToken  string `json:"accessToken"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int64  `json:"expiresIn"`
	RefreshToken string `json:"refreshToken,omitempty"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken"`
}

type SessionTokenClaims struct {
	TokenClaims

	RefreshTokenID string `json:"rti,omitempty"`
}

type LoginDTO struct {
	UserName              string `json:"userName"`
	Password              string `json:"password"`
	RefreshTokenExpiresIn int    `json:"refreshTokenExp"`
}
