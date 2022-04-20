package github

type Adapter struct {
	config Config
}

type Config struct {
	ClientID            string
	ClientSecret        string
	OAuthRedirectURI    string
	OAuthStateExpiresIn int
}

func New(config Config) Adapter {
	return Adapter{
		config: config,
	}
}
