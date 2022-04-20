package identity

type ThirdPartyTokenType int8

const (
	ThirdPartyTokenTypeGithubAccessToken ThirdPartyTokenType = 10*iota + 11
)

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)
