package error

import (
	errPayment "github.com/anddriii/kita-futsal/payment-service/constants/error/payment"
)

// ErrMapping checks if an error exists in predefined error lists
func ErrMapping(err error) bool {
	var (
		GeneralErrors = GeneralErrors
		TimerErrors   = errPayment.PaymentErrors
	)
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, TimerErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true // Error found in predefined lists
		}
	}

	return false // Error not found
}
