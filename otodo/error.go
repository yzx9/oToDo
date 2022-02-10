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

	// Resource
	ErrorNotFound ErrorCode = 30000

	// Logic
)
