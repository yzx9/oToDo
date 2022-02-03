package utils

import "fmt"

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
