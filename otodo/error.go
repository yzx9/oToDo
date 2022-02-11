package otodo

type ErrorCode int

type Error struct {
	Code    ErrorCode
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// TODO
const (
	// Request

	// Request/Auth
	ErrorUnauthorized ErrorCode = 10000
	ErrorForbidden    ErrorCode = 10001

	// Request/Limit
	ErrorRequestEntityTooLarge ErrorCode = 11000
	ErrorPreconditionFailed    ErrorCode = 11001
	ErrorPreconditionRequired  ErrorCode = 11002

	// Resource
	ErrorDuplicateID ErrorCode = 20000
	ErrorNotFound    ErrorCode = 20001

	// Logic
	ErrorAbort ErrorCode = 30000
)
