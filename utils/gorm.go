package utils

import "gorm.io/gorm"

func WrapGormErr(err error, resource string) error {
	if err != gorm.ErrRecordNotFound {
		return NewErrorWithNotFound("%v not found", resource)
	}

	return NewErrorWithUnknown("unknown error in %v: %w", resource, err)
}
