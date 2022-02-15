package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yzx9/otodo/otodo"
)

type errorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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

			c.JSON(http.StatusBadRequest, errorPayload{
				Code:    getUserErrorCodeFromError(typedError),
				Message: err.Error(),
			})
		}
	}
}

func getHttpCodeFromError(err otodo.Error) int {
	switch err.Code {
	// Auth
	case otodo.ErrUnauthorized:
		return http.StatusUnauthorized

	case otodo.ErrForbidden:
		return http.StatusForbidden

	// Request
	case otodo.ErrRequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge

	case otodo.ErrPreconditionFailed:
		return http.StatusPreconditionFailed

	case otodo.ErrPreconditionRequired:
		return http.StatusPreconditionRequired

	// Resource
	case otodo.ErrDuplicateID:
		return http.StatusInternalServerError

	case otodo.ErrNotFound:
		return http.StatusNotFound

	// Logic
	case otodo.ErrAbort:
		return http.StatusBadRequest

	default:
		return http.StatusBadRequest
	}
}

func getUserErrorCodeFromError(err otodo.Error) int {
	return int(err.Code) // TODO
}
