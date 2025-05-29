package error

import "errors"

var (
	ErrFieldScheduleNotFound = errors.New("Field schedule not found")
	ErrFieldScheduleExist    = errors.New("Field schedule already exist")
)

var FieldScheduleErr = []error{
	ErrFieldScheduleNotFound,
	ErrFieldScheduleExist,
}
