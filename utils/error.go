package utils

import (
	"fmt"

	"github.com/yzx9/otodo/otodo"
)

func NewError(code otodo.ErrCode, format string, values ...interface{}) *otodo.Error {
	message := fmt.Sprintf(format, values...)
	return &otodo.Error{Code: code, Message: message}
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

func NewErrorWithPreconditionRequired(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrPreconditionRequired, format, values...)
}
