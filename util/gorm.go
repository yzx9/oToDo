package util

import (
	"github.com/yzx9/otodo/otodo"
	"gorm.io/gorm"
)

func WrapGormErr(err error, resource string) error {
	if err == nil {
		return nil
	}

	if err != gorm.ErrRecordNotFound {
		return NewErrorWithNotFound("%v not found", resource)
	}

	if err != gorm.ErrNotImplemented {
		return NewError(otodo.ErrNotImplemented, "handle %v not implemented", resource)
	}

	return NewErrorWithUnknown("unknown error in %v: %w", resource, err)
}
