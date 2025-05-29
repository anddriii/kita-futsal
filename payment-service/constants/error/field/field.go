package error

import "errors"

var ErrFieldNotFound = errors.New("Field not found")

var FieldsErrors = []error{
	ErrFieldNotFound,
}
