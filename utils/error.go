package utils

import (
	"fmt"
	"net/http"
)

type ErrorWithHttpStatus struct {
	Code    int
	message string
}

func (err *ErrorWithHttpStatus) Error() string {
	return err.message
}

func NewErrorWithHttpStatus(code int, format string, values ...interface{}) *ErrorWithHttpStatus {
	message := fmt.Sprintf(format, values...)
	return &ErrorWithHttpStatus{Code: code, message: message}
}

func NewErrorWithForbidden(format string, values ...interface{}) *ErrorWithHttpStatus {
	return NewErrorWithHttpStatus(http.StatusForbidden, format, values...)
}

func NewErrorWithNotFound(format string, values ...interface{}) *ErrorWithHttpStatus {
	return NewErrorWithHttpStatus(http.StatusNotFound, format, values...)
}
