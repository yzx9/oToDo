package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	otodoErrors "github.com/yzx9/otodo/infrastructure/errors"
	"github.com/yzx9/otodo/model/dto"
)

func ErrorMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Next()

		if c.IsAborted() {
			err := c.Errors.Last()
			typedError := otodoErrors.Error{}
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

func getHttpCodeFromError(err otodoErrors.Error) int {
	switch err.Code {
	// Internal Error
	case otodoErrors.ErrUnknown:
		return http.StatusInternalServerError

	case otodoErrors.ErrNotImplemented:
		return http.StatusNotImplemented

	// Auth
	case otodoErrors.ErrUnauthorized:
		return http.StatusUnauthorized

	case otodoErrors.ErrForbidden:
		return http.StatusForbidden

	// Limit
	case otodoErrors.ErrPreconditionRequired:
		return http.StatusPreconditionRequired

	case otodoErrors.ErrPreconditionFailed:
		return http.StatusPreconditionFailed

	case otodoErrors.ErrRequestEntityTooLarge:
		return http.StatusRequestEntityTooLarge

	// Resource
	case otodoErrors.ErrDatabaseConnectFailed:
		return http.StatusServiceUnavailable

	case otodoErrors.ErrDataInconsistency:
		return http.StatusInternalServerError

	case otodoErrors.ErrDuplicateID:
		return http.StatusInternalServerError

	case otodoErrors.ErrNotFound:
		return http.StatusNotFound

	// Third Party
	case otodoErrors.ErrThirdPartyUnknown:
		return http.StatusBadRequest

	case otodoErrors.ErrThirdPartyUnauthorized:
		return http.StatusBadRequest

	case otodoErrors.ErrThirdPartyForbidden:
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func getUserErrorCodeFromError(err otodoErrors.Error) int {
	return int(err.Code) // TODO: some err should be hidden
}
