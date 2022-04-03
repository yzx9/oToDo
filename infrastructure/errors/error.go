package errors

// TODO: remove

type ErrCode int

type Error struct {
	Code    ErrCode
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// TODO: define err code

// Internal Error
const (
	ErrUnknown ErrCode = iota
	ErrNotImplemented
)

// Auth
const (
	ErrUnauthorized ErrCode = 10000 + iota
	ErrForbidden
)

// Limit
const (
	ErrPreconditionRequired ErrCode = 20000 + iota
	ErrPreconditionFailed
	ErrRequestEntityTooLarge
	ErrBadRequest
)

// Resource
const (
	ErrDatabaseConnectFailed ErrCode = 30000 + iota
	ErrDataInconsistency
	ErrDuplicateID
	ErrNotFound
)

// Third Party
const (
	ErrThirdPartyUnknown ErrCode = 40000 + iota
	ErrThirdPartyUnauthorized
	ErrThirdPartyForbidden
)
