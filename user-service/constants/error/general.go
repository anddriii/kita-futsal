package error

import "errors"

// General application errors
var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrToManyRequest       = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
)

// List of general errors
var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrToManyRequest,
	ErrUnauthorized,
	ErrInvalidToken,
	ErrForbidden,
}
