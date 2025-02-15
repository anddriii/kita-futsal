package error

// ErrMapping checks if an error exists in predefined error lists
func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(append(GeneralErrors[:], UserErrors[:]...)) // Merging general and user errors

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true // Error found in predefined lists
		}
	}

	return false // Error not found
}
