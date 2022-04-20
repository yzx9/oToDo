package util

import (
	"errors"

	customError "github.com/yzx9/otodo/infrastructure/errors"
	"gorm.io/gorm"
)

func WrapGormErr(err error, resource string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewErrorWithNotFound("%v not found", resource)
	}

	if errors.Is(err, gorm.ErrNotImplemented) {
		return NewError(customError.ErrNotImplemented, "handle %v not implemented", resource)
	}

	return NewErrorWithUnknown("unknown error in %v: %w", resource, err)
}
