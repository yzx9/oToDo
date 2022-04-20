package config

var GitHub = Config{}

type Config struct {
	ClientID            string
	ClientSecret        string
	OAuthRedirectURI    string
	OAuthStateExpiresIn int
}
