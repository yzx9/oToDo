package identity

import (
	"fmt"
)

var InvalidCredential = fmt.Errorf("invalid credential.")
var UserNameDuplicate = fmt.Errorf("user name has been used.")
var UserNameTooShort = fmt.Errorf("user name too short.")
var InvalidPassword = fmt.Errorf("invalid password.")

type Error struct {
	err error
}

func (e Error) Error() string { return "unknown error." }
func (e Error) Unwrap() error { return e.err }
