package identity

import (
	"fmt"
)

var InvalidCredential = fmt.Errorf("invalid credential.")
var UserNameDuplicate = fmt.Errorf("duplicate.")
var UserNameTooShort = fmt.Errorf("user name too short.")
var PasswordTooShort = fmt.Errorf("password too short.")

type InnerError struct {
	err error
}

func (InnerError) Error() string {
	return "identity domain error."
}

func (e InnerError) Unwrap() error {
	return e.err
}

func newErr(err error) InnerError {
	return InnerError{err: err}
}
