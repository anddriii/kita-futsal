package error

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

//When the front end makes a request to the service, if the request is invalid, it will be handled by this common error.

type ValidateResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{}

func ErrValidationResponse(err error) (validationResponse []ValidateResponse) {
	var fieldErrors validator.ValidationErrors
	// Check if the error matches the validator's error type
	if errors.As(err, &fieldErrors) {
		for _, err := range fieldErrors {
			switch err.Tag() {
			case "required":
				validationResponse = append(validationResponse, ValidateResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is required", err.Field()),
				})
			case "email":
				validationResponse = append(validationResponse, ValidateResponse{
					Field:   err.Field(),
					Message: fmt.Sprintf("%s is not valid email address", err.Field()),
				})
			default:
				// Check if there is a predefined error message for this validation tag
				errValidator, ok := ErrValidator[err.Tag()]
				if ok {
					count := strings.Count(errValidator, "%s")
					if count == 1 {
						validationResponse = append(validationResponse, ValidateResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field()),
						})
					} else {
						validationResponse = append(validationResponse, ValidateResponse{
							Field:   err.Field(),
							Message: fmt.Sprintf(errValidator, err.Field(), err.Param()),
						})
					}
				} else {
					// Default error message if no predefined message exists
					validationResponse = append(validationResponse, ValidateResponse{
						Field:   err.Field(),
						Message: fmt.Sprintf("something wrong on %s; %s", err.Field(), err.Tag()),
					})
				}
			}
		}
	}
	return validationResponse
}

func WrapError(err error) error {
	logrus.Errorf("error %v", err)
	return err
}
