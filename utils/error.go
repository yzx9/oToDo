package utils

import (
	"fmt"

	"github.com/yzx9/otodo/otodo"
)

func NewError(code otodo.ErrorCode, format string, values ...interface{}) *otodo.Error {
	message := fmt.Sprintf(format, values...)
	return &otodo.Error{Code: code, Message: message}
}

func NewErrorWithForbidden(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrorForbidden, format, values...)
}

func NewErrorWithNotFound(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrorNotFound, format, values...)
}

func NewErrorWithPreconditionFailed(format string, values ...interface{}) *otodo.Error {
	return NewError(otodo.ErrorPreconditionFailed, format, values...)
}
