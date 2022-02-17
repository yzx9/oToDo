package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/otodo"
)

func ErrorMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Next()

		if c.IsAborted() {
			err := c.Errors.Last()
			typedError := otodo.Error{}
			if errors.As(err, &typedError) {
				code := getHttpCodeFromError(typedError)
				c.AbortWithError(code, typedError)
				return
			}

			c.JSON(http.StatusBadRequest, dto.ErrorDTO{
				Code:    getUserErrorCodeFromError(typedError),
				Message: err.Error(),
			})
		}
	}
}

func getHttpCodeFromError(err otodo.Error) int {
	switch err.Code {
	// Internal Error
	case otodo.ErrUnknown:
		return http.StatusInternalServerError

	case otodo.ErrNotImplemented:
		return http.StatusNotImplemented

	// Auth
	case otodo.ErrUnauthorized:
		return http.StatusUnauthorized

	case otodo.ErrForbidden:
		return http.StatusForbidden

	// Limit
	case otodo.ErrPreconditionRequired:
		return http.StatusPreconditionRequired

	case otodo.ErrPreconditionFailed:
		return http.StatusPreconditionFailed

	case otodo.ErrRequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge

	// Resource
	case otodo.ErrDatabaseConnectFailed:
		return http.StatusServiceUnavailable

	case otodo.ErrDataInconsistency:
	case otodo.ErrDuplicateID:
		return http.StatusInternalServerError

	case otodo.ErrNotFound:
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func getUserErrorCodeFromError(err otodo.Error) int {
	return int(err.Code) // TODO: some err should be hidden
}
