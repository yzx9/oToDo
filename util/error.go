package util

import (
	"fmt"

	"github.com/yzx9/otodo/otodo"
)

func NewError(code otodo.ErrCode, format string, values ...interface{}) *otodo.Error {
	err := fmt.Errorf(format, values...)
	return &otodo.Error{Code: code, Message: err.Error()}
}

func NewErrorWithBadRequest(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrBadRequest, format, values...)
}

func NewErrorWithForbidden(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrForbidden, format, values...)
}

func NewErrorWithNotFound(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrNotFound, format, values...)
}

func NewErrorWithPreconditionFailed(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrPreconditionFailed, format, values...)
}

func NewErrorWithUnknown(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrUnknown, format, values...)
}
