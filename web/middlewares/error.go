package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

			c.AbortWithError(http.StatusBadRequest, err)
		}
	}
}

func getHttpCodeFromError(err otodo.Error) int {
	switch err.Code {
	// Auth
	case otodo.ErrorUnauthorized:
		return http.StatusUnauthorized

	case otodo.ErrorForbidden:
		return http.StatusForbidden

	// Request
	case otodo.ErrorRequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge

	case otodo.ErrorPreconditionFailed:
		return http.StatusPreconditionFailed

	case otodo.ErrorPreconditionRequired:
		return http.StatusPreconditionRequired

	// Resource
	case otodo.ErrorDuplicateID:
		return http.StatusInternalServerError

	case otodo.ErrorNotFound:
		return http.StatusNotFound

	// Logic
	case otodo.ErrorAbort:
		return http.StatusBadRequest

	default:
		return http.StatusBadRequest
	}
}
