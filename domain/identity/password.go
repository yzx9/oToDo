package identity

import (
	"bytes"
	"crypto/sha256"
)

// value object
type password struct {
	bytes []byte
}

func NewPasswordByBytes(bytes []byte) password { return password{bytes: bytes} }
func NewPassword(pwd string) (password, error) {
	if pwd == "" {
		return password{
			bytes: nil,
		}, nil
	}

	if len(pwd) < 6 {
		return password{}, InvalidPassword
	}

	if len(pwd) > 20 {
		return password{}, InvalidPassword
	}

	return password{
		bytes: cryptoPassword(pwd),
	}, nil
}

func (pwd password) Bytes() []byte { return pwd.bytes }
func (pwd password) Empty() bool   { return pwd.bytes == nil }

func (pwd password) Equals(aPassword string) bool {
	if pwd.Empty() {
		return false
	}

	crypto := cryptoPassword(aPassword)
	return bytes.Equal(pwd.bytes, crypto)
}

func cryptoPassword(password string) []byte {
	pwd := sha256.Sum256(append([]byte(password), Conf.PasswordNonce...))
	return pwd[:]
}
