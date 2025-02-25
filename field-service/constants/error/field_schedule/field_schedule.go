package error

import "errors"

var (
	ErrFieldNotScheduleFound = errors.New("Field schedule not found")
	ErrFieldScheduleExist    = errors.New("Field schedule already exist")
)

var FieldScheduleErr = []error{
	ErrFieldNotScheduleFound,
	ErrFieldScheduleExist,
}
