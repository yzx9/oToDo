package user

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)
