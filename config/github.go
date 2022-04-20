package config

var GitHub = github{}

type github struct {
	ClientID            string
	ClientSecret        string
	OAuthRedirectURI    string
	OAuthStateExpiresIn int
}
