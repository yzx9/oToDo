package config

var Session = session{}

type session struct {
	AccessTokenExpiresIn         int // second
	RefreshTokenExpiresInDefault int // second
	RefreshTokenExpiresInMax     int // second
	RefreshTokenExpiresInOAuth   int // second
	AccessTokenRefreshThreshold  int // second
}
