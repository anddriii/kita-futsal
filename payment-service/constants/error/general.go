package error

import "errors"

// General application errors
var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrToManyRequest       = errors.New("too many requests")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrSizeTooBig          = errors.New("size too big")
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
