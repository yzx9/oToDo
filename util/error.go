package util

import (
	"fmt"

	"github.com/yzx9/otodo/infrastructure/errors"
)

func NewError(code errors.ErrCode, format string, values ...interface{}) *errors.Error {
	err := fmt.Errorf(format, values...)
	return &errors.Error{Code: code, Message: err.Error()}
}

func NewErrorWithBadRequest(format string, values ...interface{}) *errors.Error {
	return NewError(errors.ErrBadRequest, format, values...)
}

func NewErrorWithForbidden(format string, values ...interface{}) *errors.Error {
	return NewError(errors.ErrForbidden, format, values...)
}

func NewErrorWithNotFound(format string, values ...interface{}) *errors.Error {
	return NewError(errors.ErrNotFound, format, values...)
}

func NewErrorWithPreconditionFailed(format string, values ...interface{}) *errors.Error {
	return NewError(errors.ErrPreconditionFailed, format, values...)
}

func NewErrorWithUnknown(format string, values ...interface{}) *errors.Error {
	return NewError(errors.ErrUnknown, format, values...)
}
