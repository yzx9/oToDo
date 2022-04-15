package session

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
)
