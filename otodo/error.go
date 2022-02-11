package otodo

type ErrCode int

type Error struct {
	Code    ErrCode
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// TODO
const (
	// Request

	// Request/Auth
	ErrUnauthorized ErrCode = 10000
	ErrForbidden    ErrCode = 10001

	// Request/Limit
	ErrRequestEntityTooLarge ErrCode = 11000
	ErrPreconditionFailed    ErrCode = 11001
	ErrPreconditionRequired  ErrCode = 11002

	// Resource
	ErrDuplicateID ErrCode = 20000
	ErrNotFound    ErrCode = 20001

	// Logic
	ErrAbort ErrCode = 30000
)
