package error

import "errors"

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrExpireArInvalid = errors.New("expire is invalid, must be greater than current time")
)

var PaymentErrors = []error{
	ErrExpireArInvalid,
	ErrPaymentNotFound,
}
