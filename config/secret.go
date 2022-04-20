package config

var Secret = secret{}

type secret struct {
	TokenIssuer     string
	TokenHmacSecret []byte
	PasswordNonce   []byte
}
